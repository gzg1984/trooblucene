package main

import (
	"fmt"
	"os"

	std "github.com/gzg1984/golucene/analysis/standard"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/gzg1984/golucene/core/codec/lucene71"
	"github.com/gzg1984/golucene/core/document"
	"github.com/gzg1984/golucene/core/index"
	"github.com/gzg1984/golucene/core/search"
	"github.com/gzg1984/golucene/core/store"
	"github.com/gzg1984/golucene/core/util"
	troobconfig "github.com/gzg1984/troobconfig"

	"github.com/gzg1984/golucene/queryparser/classic"
)

func main() {

	spdkIndexPath := troobconfig.GetIndexPath("spdk")
	fmt.Printf("spdkIndexPath is %v\n", spdkIndexPath)

	util.SetDefaultInfoStream(util.NewPrintStreamInfoStream(os.Stdout))
	//util.SetDefaultInfoStream(util.NO_OUTPUT)

	index.DefaultSimilarity = func() index.Similarity {
		return search.NewDefaultSimilarity()
	}

	//directory, _ := store.OpenFSDirectory("test_index")

	//createTestIndex(directory)
	//searchTestIndex(directory)

	sourceIndexDirectory, _ := store.OpenFSDirectory(spdkIndexPath)
	searchSource(sourceIndexDirectory, "LICENSE")
}

func simpleSearch(directory store.Directory, query string) {
	fmt.Printf("=====Enter simpleSearch\n")
	reader, _ := index.OpenDirectoryReader(directory)
	searcher := search.NewIndexSearcher(reader)

	/*TermQuery*/
	q := search.NewTermQuery(index.NewTerm("foo", query))
	res, _ := searcher.Search(q, nil, 1000)
	fmt.Printf("=====NewTermQuery Found %v hit(s) for %v.\n",
		res.TotalHits, query)
	for _, hit := range res.ScoreDocs {
		fmt.Printf("Doc %v score: %v\n", hit.Doc, hit.Score)
		doc, _ := reader.Document(hit.Doc)
		fmt.Printf("foo -> %v\n", doc.Get("foo"))
	}
}

func searchSource(directory store.Directory, key string) {
	fmt.Printf("=====Enter searchSource\n")

	/* multi field*/
	//MultiFieldQueryParser(fields, lxranalyzer).parse(lxrqueryString)
	reader, _ := index.OpenDirectoryReader(directory)
	searcher := search.NewIndexSearcher(reader)
	stopWords := make(map[string]bool)
	analyzer := std.NewStandardAnalyzerWithStopWords(stopWords)

	qp := classic.NewQueryParser(util.VERSION_LATEST, "content", analyzer)
	multq, _ := qp.Parse(key)
	res, _ := searcher.Search(multq, nil, 1000)
	fmt.Printf("===========searchSource Found %v hit(s).\n", res.TotalHits)
	for _, hit := range res.ScoreDocs {
		fmt.Printf("Doc %v score: %v\n", hit.Doc, hit.Score)
		doc, _ := reader.Document(hit.Doc)
		fmt.Printf("content -> %v\n",
			doc.Get("filePath")+"\t"+doc.Get("fileName")+"\t"+doc.Get("projectId"))
	}
}

func searchTestIndex(directory store.Directory) {
	simpleSearch(directory, "bar")
	simpleSearch(directory, "bar1")
	simpleSearch(directory, "bar2")
	simpleSearch(directory, "b")

	/* multi field*/
	//MultiFieldQueryParser(fields, lxranalyzer).parse(lxrqueryString)
	reader, _ := index.OpenDirectoryReader(directory)
	searcher := search.NewIndexSearcher(reader)
	analyzer := std.NewStandardAnalyzer()

	qp := classic.NewQueryParser(util.VERSION_LATEST, "foo", analyzer)
	multq, _ := qp.Parse("bar")
	res, _ := searcher.Search(multq, nil, 1000)
	fmt.Printf("NewQueryParser Found %v hit(s).\n", res.TotalHits)
	for _, hit := range res.ScoreDocs {
		fmt.Printf("Doc %v score: %v\n", hit.Doc, hit.Score)
		doc, _ := reader.Document(hit.Doc)
		fmt.Printf("foo -> %v\n", doc.Get("foo"))
	}
}

func createTestIndex(directory store.Directory) {
	analyzer := std.NewStandardAnalyzer()

	conf := index.NewIndexWriterConfig(util.VERSION_LATEST, analyzer)
	writer, _ := index.NewIndexWriter(directory, conf)

	d := document.NewDocument()
	d.Add(document.NewTextFieldFromString("foo", "bar is good", document.STORE_YES))
	writer.AddDocument(d.Fields())

	d2 := document.NewDocument()
	d2.Add(document.NewTextFieldFromString("foo", "bar is bad", document.STORE_YES))
	writer.AddDocument(d2.Fields())
	writer.Close() // ensure index is written

}
