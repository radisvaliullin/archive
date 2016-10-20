// is easy sample of analyzes parser from "https://www.invitro.ru/analizes/for-doctors/"

package main

import (
    "fmt"
    "strings"
    "os"

    "golang.org/x/net/html"
    "net/http"
    "golang.org/x/net/html/charset"
    //"time"
    "time"
)


// Parsing page urls
const domain_url string = "https://www.invitro.ru"
const analyzes_url string = "/analizes/for-doctors/"
const base_url string = domain_url + analyzes_url
// details tags
const detailDesc string = "link_88"
const detailPrep string = "link_84"
const detailIndic string = "link_81"
const detailInterp string = "link_82"


var test_cnt int = 0


func main() {

    //
    // getting the base_url page's html.
    //
    // get http.Get respond
    resp, err := http.Get(base_url)
    if err != nil {
        fmt.Println("Error http.Get:", err)
        return
    }
    defer resp.Body.Close()
    // convert to UTF-8
    html_utf8_rd, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
    if err != nil {
        fmt.Println("Error converting win-1251 to utf-8:", err)
        return
    }

    // html from saved file (for testing).
    //html_utf8_rd, err := os.Open("test_base_page.html")
    //if err != nil {
    //    fmt.Println("Error open file:", err)
    //    return
    //}
    //defer html_utf8_rd.Close()

    // html parsing, get html nodes tree
    html_tree, err :=	html.Parse(html_utf8_rd)
    if err != nil {
        fmt.Println("Error html.Parse:", err)
        return
    }

    // get analyzes catalog nodes tree
    analyzes_tree := getAnalyzesCatalogTree(html_tree)

    // chan for writing analyzes to DB
    analyse_chan := make(chan Analyse)

    // run to DB write goruotine
    go writeToDB(analyse_chan)

    // analyzes catalog parsing
    parseAnalyzesTree(&ParserState{}, analyzes_tree, analyse_chan)

    time.Sleep(time.Second * 10)

}


func getAnalyzesCatalogTree(n *html.Node) *html.Node {
    var target_node *html.Node
    if n.Type == html.ElementNode && n.Data == "div" {
        for _, a := range n.Attr {
            if a.Key == "id" && a.Val == "catalog-section-analiz" {
                target_node = n
            }
        }
    }
    if target_node == nil {
        for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
            target_node = getAnalyzesCatalogTree(ch)
            if target_node != nil { break }
        }
    }
    return target_node
}


func getDetailsNode(n *html.Node) *html.Node {
    var target_node *html.Node
    if n.Type == html.ElementNode && n.Data == "script" {
        for _, a := range n.Attr {
            if a.Key == "language" && a.Val == "JavaScript" && strings.HasPrefix(strings.TrimSpace(n.FirstChild.Data), "var arTexts=") {
                target_node = n
            }
        }
    }
    if target_node == nil {
        for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
            target_node = getDetailsNode(ch)
            if target_node != nil { break }
        }
    }
    return target_node
}


type Analyse struct {
    Type string
    Subtype string
    Kind string
    DetailUrl string
    DetailDesc string
    DetailPrep string
    DetailIndic string
    DetailInterp string
}


type ParserState struct {
    level int
    lv_1_name string
    lv_2_name string
}


func parseAnalyzesTree(state *ParserState, n *html.Node, analyse_chan chan Analyse) {
    // analyse parsing limit
    if test_cnt > 20 { return }

    if n.Type == html.ElementNode {
        // set level in analyzes catalog
        if n.Data == "td" {
            for _, a := range n.Attr {
                if a.Key == "class" {
                    if a.Val == "name sec_1" {
                        state.level = 1
                    } else if a.Val == "name sec_2" {
                        state.level = 2
                    } else if a.Val == "name sec_3" {
                        state.level = 3
                    }else if a.Val == "name sec_4" {
                        state.level = 4
                    }
                }
            }
        // get 1 & 2 level names
        } else if n.Data == "span" && (state.level == 1 || state.level == 2) {
        // } else if n.Data == "span" && level == 1 {
            for _, a := range n.Attr {
                if a.Key == "class" && (a.Val == "name_gr") {
                    if state.level == 1 {
                        state.lv_1_name = n.FirstChild.Data
                        state.lv_2_name = ""
                    } else if state.level == 2 {
                        state.lv_2_name = n.FirstChild.Data
                    }
                }
            }
        // get analyse kind
        // }
        } else if n.Data == "a" {
            for _, a := range n.Attr {
                if a.Key == "href" && strings.HasPrefix(a.Val, "/") {
                    analyse := Analyse{
                        Type:state.lv_1_name,
                        Subtype:state.lv_2_name,
                        Kind:n.FirstChild.Data,
                        DetailUrl:a.Val,
                    }

                    // get analyse detail in goroutine
                    go getAnalyseDetailByURL(analyse, analyse_chan)
                    test_cnt += 1
                }
            }
        }
    }
    for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
        parseAnalyzesTree(state, ch, analyse_chan)
    }
}


func getDetailsNodeText(lv int, out_tree[]string, n *html.Node) []string {
    lv += 2
    if n.Type == html.ElementNode {
        if n.Data == "br" {
            out_tree = append(out_tree, "\n")
        } else if n.Data == "div" && len(out_tree) != 0 && !strings.Contains(out_tree[len(out_tree)-1], "\n") {
            out_tree = append(out_tree, "\n")
        } else {
            for _, a := range n.Attr {
                if a.Key == "href" {
                    out_tree = append(out_tree, " [" + a.Val + "] ")
                }
            }
        }
    } else if n.Type == html.TextNode {
        if strings.TrimSpace(n.Data) != "" {
            out_tree = append(out_tree, n.Data)
        }
    }
    for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
        out_tree = getDetailsNodeText(lv, out_tree, ch)
    }
    return out_tree
}


func writeToDB(analyse_chan chan Analyse) {
    file, err := os.Create("analyzes_hierarchy.txt")
    if err != nil {
        fmt.Println("Analyzes hierarchy.", err)
        return
    }
    defer file.Close()

    for a := range analyse_chan {

        //a = getAnalyseDetailByURL(a)

        file.WriteString("--" + a.Type + "\n")
        file.WriteString("----" + a.Subtype + "\n")
        file.WriteString("------" + a.Kind + "\n")
        file.WriteString("--------" + a.DetailDesc + "\n")
        file.WriteString("--------" + a.DetailPrep + "\n")
        file.WriteString("--------" + a.DetailIndic + "\n")
        file.WriteString("--------" + a.DetailInterp + "\n")
    }
}


func getAnalyseDetailByURL(analyse Analyse, analyse_chan chan Analyse) {

    // get analyse kind's detail information
    //
    // html parsing, get html nodes tree for detail
    //
    // invalid url
    if !strings.HasPrefix(analyse.DetailUrl, analyzes_url) {return}

    // get http.Get respond
    resp, err := http.Get(domain_url + analyse.DetailUrl)
    if err != nil {
        fmt.Println("Error http.Get:", err)
        return
    }
    defer resp.Body.Close()
    // convert to UTF-8
    details_html_rd, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
    if err != nil {
        fmt.Println("Error converting win-1251 to utf-8:", err)
        return
    }
    // for testing get from file.
    //details_html_rd, _ := os.Open("test_156_28932.html")
    //defer details_html_rd.Close()

    details_html_tree, _ :=	html.Parse(details_html_rd)
    details_node := getDetailsNode(details_html_tree)
    if details_node == nil {return}
    details_subnodes_raw := details_node.FirstChild.Data

    details_subnodes_htmls := getDetailsSubnodesHTMLs(details_subnodes_raw)

    for det, d_html := range details_subnodes_htmls {
        detail, _ := html.Parse(strings.NewReader(d_html))
        detail_text := strings.Join(getDetailsNodeText(0, nil, detail), "")
        detail_text = strings.TrimSpace(detail_text)
        if det == detailDesc {
            analyse.DetailDesc = detail_text
        } else if det == detailPrep {
            analyse.DetailPrep = detail_text
        } else if det == detailIndic {
            analyse.DetailIndic = detail_text
        } else if det == detailInterp {
            analyse.DetailInterp = detail_text
        }
    }
    // set analyse to chan
    analyse_chan <- analyse
}


func getDetailsSubnodesHTMLs(details_raw string) map[string]string {
    subnodes := make(map[string]string)
    sections := []string{detailInterp, detailIndic, detailPrep, detailDesc, }
    details_raw = strings.TrimSpace(details_raw)
    details_raw = strings.TrimSuffix(details_raw, "'};")
    for _, sec := range sections {
        // dividet element
        div := "','" + sec + "':'"
        if sec == detailDesc {
             div = "'" + sec + "':'"
        }
        if strings.Contains(details_raw, div) {
            splits := strings.Split(details_raw, div)
            // delete/replase symbols
            subnodes[sec] = strings.Replace(strings.Replace(splits[1], "\\n", "", -1), "\\", "", -1)
            details_raw = splits[0]
        } else {
            subnodes[sec] = ""
        }
    }
    return subnodes
}
