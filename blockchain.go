package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"../uri"
	"../uri/handlers"
	"log"
	"net/http"
	"os"
)
/*Block Struct */
type Block struct{
	Header Header `json:""`
	

}

/* Header Struct *
* Used inside Block Struct */

type Header struct{
	Height     int32  `json:"height"`
	Timestamp  int64  `json:"timestamp"`
	Hash       string `json:"hash"`
	ParentHash string `json:"parentHash"`
	Size       int32  `json:"size"`
	Value      string `json:"value"`
	Nonce	   int    `json:"nonce"`
}

//BlockChain Struct
type BlockChain struct{
	Chain map[int32][]Block
	Length int32
}

//Block functions

/* This function takes arguments(such as height, parentHash, and value) *
* and forms a block. This is a method of the block struct.*/

func (block *Block) Initial(height int32, parentHash string, value string ){
	header  := Header{};
	header.Height = height;
	header.Timestamp = GetTimeStamp();
	header.ParentHash = parentHash;
	header.Size = 32;
	header.Value = value;
	hashStr := string(header.Height)+ string(header.Timestamp)+ header.ParentHash + string(header.Size)+ header.Value;
	header.Hash = generateHash(hashStr);
	block.Header = header;
}

/*This function takes a string that represents the JSON value of a block *
* as an input, and decodes the input string back to a block instance. *
* Argument: a string of JSON format *
* Return value: a block instance */

func DecodeFromJson(jsonString string) Block{
	var block Block;
	_ = json.Unmarshal([]byte(jsonString), &block);
	return block;
}

/*This function takes an array of JSON strings and converts them into an *
* array of blocks. This is used for the Blockchain function DecodeFromJson() *
*instead of the Block Decode from Json because now we have to deal with an Array */

func DecodeFromJsonArray(jsonString string) []Block{
	var block[] Block;
	_ = json.Unmarshal([]byte(jsonString), &block);
	return block;
}

/*This function encodes a block instance into a JSON format string. *
* Argument: a block or you may define this as a method of the block struct *
* Return value: a string of JSON format*/

func (block *Block) EncodeToJson()string{
	jsonBlock,_ := json.Marshal(block);
	return string(jsonBlock);
}

//BlockChain Functions

/*This function takes a height as the argument, returns the list of *
*blocks stored in that height or None if the height doesn't exist. *
* Argument: int32 *
* Return type: []Block */

func (chain *BlockChain) Get(height int32) []Block{

	if(chain.Chain[height] == nil){
		return nil;
	}
	return chain.Chain[height];
}

/*This function generates a nonce that produces a hash with the *
* correct difficulty *
* Return type: int *
 */

func (block *Block) GetNonce() int{

	n := rand.Intn(1000)
	var nonce int;
	for i:= 0; ; i++{
		newHash := block.Header.ParentHash + strconv.Itoa(n+i) + block.Header.Value
		testHash := generateNonceHash(newHash)
		pass := CheckHash(5, testHash)
		if pass == true{
			nonce = n + i;
			block.Header.Nonce = nonce
			return nonce
		}
	}


}

/* This function takes the desired difficulty (# of zeros needed) and *
*  checks if the given hash starts with the desired # of 0's *
*  Arguments: int, string  *
* Return type: int (1 on success and 0 on failure)  *
 */

func CheckHash(difficulty int, hash string) bool{
	if difficulty > 256{ //difficulty is greater than hash elements (for sha 256)
		return false
	}
	prefix := strings.Repeat("0", difficulty)
	//var zeroHash  string
/*
	for i:= 0; i < difficulty; i++{
		zeroHash = zeroHash + "0"
	}*/
	//compareHash := string(hash[0:difficulty])
	return strings.HasPrefix(hash, prefix)

}
/* This function takes a block as the argument, use its height to find *
* the corresponding list in blockchain's Chain map. If the list has already *
* contained that block's hash, ignore it because we don't store duplicate blocks; *
* if not, insert the block into the list. *
* Argument: block */

func (chain *BlockChain) Insert(block Block){

	hashCheck := block.Header.ParentHash + strconv.Itoa(block.Header.Nonce) + block.Header.Value
	testHashCheck := generateNonceHash(hashCheck)


	/*checks if block matches the desired difficulty */

	if CheckHash(5, testHashCheck) == false {
		return
	}


	height := block.Header.Height

	/* checks if block is already in the chain */

	for i := 0; i < len(chain.Chain[height]); i++{
		if block.Header.Hash == chain.Chain[height][i].Header.Hash{
			return
		}

	}

	/*checks if parent block is in chain */

	if height > 0 {
		for i := 0; i < len(chain.Chain[height-1]); i++ {
			if block.Header.ParentHash == chain.Chain[height-1][i].Header.Hash{
				return   //parent is not in the chain
			}
		}
	}

	if height > chain.Length{
		chain.Length = height;
	}
	chain.Chain[height] = append(chain.Chain[height], block);


}

/*This function iterates over all the blocks, generate blocks' JSON String *
* by the function you implemented for a single Block, and return the list of *
* those JSON Strings. *
* Return type: string */

func (chain *BlockChain) EncodeToJson() string{

	var jsonS string;

	jsonS = `[`;
	var i int32;
	for i = 1; i <= chain.Length; i++{

		for j:= 0; j < len(chain.Chain[i]); j++  {

			if i == chain.Length && j == len(chain.Chain[i])-1{
				jsonS += chain.Chain[i][j].EncodeToJson();
			}else{
				jsonS += chain.Chain[i][j].EncodeToJson() + `,`;
			}
		}
		fmt.Println();

	}
	jsonS = jsonS + `]`;


	return jsonS;
}

/*This function is called upon a blockchain instance. It takes a blockchain JSON *
* string as input, decodes the JSON string back to a list of block JSON strings, *
* decodes each block JSON string back to a block instance, and inserts every block *
* into the blockchain. *
* Argument: self, string */
func  (chain *BlockChain) DecodeFromJson(jsonString string) {

	var block[] Block;
	block = DecodeFromJsonArray(jsonString);


	for i:= 0; i < len(block); i++{
		chain.Insert(block[i]);
	}
}

/*Helper function that creates Hashed Values*/
func generateHash(value string) string{
	hash := sha256.New()
	hash.Write([]byte(value))
	codeInBytes := hash.Sum(nil)
	codedString := hex.EncodeToString(codeInBytes[:])
	value2, err := json.Marshal(codedString)
	if err != nil{
		fmt.Print("error: %v", err)
	}
	if err != nil{
		log.Fatal(err)
	}
	return string(value2)
}

/*Helper function that creates Hashed  for PoW Values*/
func generateNonceHash(value string) string{
	hash := sha256.New()
	hash.Write([]byte(value))
	codeInBytes := hash.Sum(nil)
	codedString := hex.EncodeToString(codeInBytes[:])

	return codedString
}

/*Helper function to get timestamp*/

func GetTimeStamp() int64{
	currentTime := time.Now()
	seconds := currentTime.Unix();
	return seconds;
}

//Main
func main(){
	rand.Seed(time.Now().UnixNano())

	router := uri.NewRouter()

	var port string
	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		port = "6689"
	}

	handlers.InitSelfAddress(port)

	log.Fatal(http.ListenAndServe(":"+port, router))

	block := Block{}
	block.Initial(1, "0(genesis block)",  "");
	block.GetNonce()
	//block2:= Block{};
	//block2.Initial(2, block.Header.Hash, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08");
	block3:= Block{};
	block3.Initial(2, block.Header.Hash, "9f86d081884c7d861bce3b28f33eec1be758a211b2b0b822cd15d6c15b0a08");
	block4 :=Block{};
	block4.Initial(3, block3.Header.Hash, "60303ae22b998861bce3b28f33eec1be758a211b2b0b822cd15d6c15b0f00a08");
	block5 :=Block{};
	block5.Initial(4, block4.Header.Hash, "50303ae22b998861bce3b28f3stew1b2b0b822cd15d6c15b0f00a08");
	chain := BlockChain{};
	chain.Chain = make(map[int32][]Block);
	chain.Length = 0;
	chain.Insert(block);
	//chain.Insert(block2);
	chain.Insert(block3);
	chain.Insert(block4);
	chain.Insert(block5);
	var str string;
	str  = chain.EncodeToJson();

	chain2 := BlockChain{
		Chain:  make(map[int32][]Block),
		Length: 0,
	}
	chain2.DecodeFromJson(str);
	//fmt.Println(chain.Length);
	fmt.Println(chain);

	fmt.Println(str);




}




