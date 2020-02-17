package main

import (
	"github.com/kuso/japanese-word-extractor/server"
)

func main() {
	/*
	dict, err := jlpt_tool.NewJLPTDictionary()
	if err != nil {
		log.Fatal(err)
	}

	//var str = "寿司が食べたい。"
	//var str = "電車でんしゃの中なかで体からだを触さわられたりする「痴漢ちかん」の被害ひがいをなくすため、ＪＲ東日本ひがしにほんはスマートフォンのアプリを使つかったシステムを作つくりました。"
	var str = `レオパレス21は27日、大株主の投資会社レノ（東京・渋谷）が開催を求めていた臨時株主総会を2月27日に開くと発表した。臨時株主総会では、社外取締役を2人追加する会社議案を提案し、企業統治改善で株主の支持を求める。レオパレスは取締役10人全員の解任などを求めるレノの株主提案に反対しており、経営権を巡る争いは臨時株主総会で山場をむかえる。
	レオパレスは27日開催した取締役会で、東洋シヤッター元社長の藤田和育氏とパナソニックホームズ元上席主幹の中村裕氏を社外取締役候補として決めた。レオパレスの取締役は現在、社内5人、社外5人の計10人。2人が新たに加われば、社外取締役が過半となる。これまでは6月の定時株主総会で社外取締役を過半数とする方針だったが、前倒しする。
	レノは著名投資家の村上世彰氏が関与するファンド。レオパレスの宮尾文也社長ら取締役10人全員の解任と、レノが推薦する取締役3人の選任を求めている。レオパレスは、レノ側の提案に「自己の利益を追求する目的」として反対している。　会社提案と株主提案が可決されるには総会に出席する株主の過半数の賛成が必要となる。レオパレスはレノが14.46%（共同保有含む）、英運用会社オデイ・アセット・マネジメントが14.34%、国内運用会社のアルデシアインベストメントが16.1%の株式を持ち、3社で発行済み株式数の約45%を保有する。オデイとアルデシアは態度を明らかにしていない。`
	tokens := jlpt_tool.GetTokens(str, dict)
	//jlpt_tool.PrintTokens(tokens)
	html := jlpt_tool.OutputHTML(tokens)
	log.Println(html)
	log.Println(len(tokens))
	 */

	/*
	t := tokenizer.New()
	//tokens := t.Tokenize("寿司が食べたい。") // t.Analyze("寿司が食べたい。", tokenizer.Normal)
	tokens := t.Tokenize("寿司が食べたい。")

	for _, token := range tokens {
		if token.Class == tokenizer.DUMMY {
			// BOS: Begin Of Sentence, EOS: End Of Sentence.
			fmt.Printf("%s\n", token.Surface)
			continue
		}
		if token.Surface == "\n" {
			fmt.Printf("%s", "linebreak")
		}
		fmt.Printf("%s", token.Surface)
	}
	os.Exit(0)
	 */

	queueName := "test_jlpt_queue"
	server := server.NewServer()
	server.SetupRouter()
	server.SetupMQ("test_jlpt_service", queueName)

	go server.HttpServer.ListenAndServe()
	server.GracefulShutdown(3000)
}
