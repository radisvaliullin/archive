// is easy sample of analyzes parser from "https://www.invitro.ru/analizes/for-doctors/"

package main

import (
    "fmt"
    "strings"

    "golang.org/x/net/html"
    "net/http"
    "golang.org/x/net/html/charset"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"

    "sync"
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

// node parsing limits, for testing
var test_cnt int = 0
// sync.WaitGroup
var parser_wg sync.WaitGroup
var write_db_wg sync.WaitGroup


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
    //html_utf8_rd, err := os.Open("test_html/test_base_page.html")
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

    // chan for writing to DB analyzes
    analyse_chan := make(chan Analyse)

    // run to DB write goruotine
    write_db_wg.Add(1)
    go writeToDB(analyse_chan)

    // analyzes catalog parsing
    parseAnalyzesTree(&ParserState{}, analyzes_tree, analyse_chan)

    //time.Sleep(time.Second * 5)
    parser_wg.Wait()
    close(analyse_chan)
    write_db_wg.Wait()
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
    if test_cnt >= 50 { return }

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
                    parser_wg.Add(1)
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


func getDetailsNodeText(out_tree[]string, n *html.Node) []string {
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
        out_tree = getDetailsNodeText(out_tree, ch)
    }
    return out_tree
}


func writeToDB(analyse_chan chan Analyse) {
    defer write_db_wg.Done()

    // open DB
    db, err := sql.Open("sqlite3", "./analyzes.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()
    // create tables
    sql_request := `
        CREATE TABLE 'analyse_type'
            (
                'id' INTEGER PRIMARY KEY AUTOINCREMENT,
                'type_name' TEXT NULL
            );
        CREATE TABLE 'analyse_subtype'
            (
                'id' INTEGER PRIMARY KEY AUTOINCREMENT,
                'subtype_name' TEXT NULL,
                'type' INTEGER NULL,
                constraint fk_analysesubtype_analysetype foreign key (type) references analyse_type (id)
            );
        CREATE TABLE 'analyse'
            (
                'id' INTEGER PRIMARY KEY AUTOINCREMENT,
                'type' INTEGER NULL,
                'subtype' INTEGER NULL,
                'analyse_kind' TEXT NULL,
                'decribe' TEXT NULL,
                'prepare' TEXT NULL,
                'indications' TEXT NULL,
                'interpritation' TEXT NULL,
                constraint fk_analyse_analysetype foreign key (type) references analyse_type (id),
                constraint fk_analyse_analysesubtype foreign key (subtype) references analyse_subtype (id)
            );

    `
    _, err = db.Exec(sql_request)
    if err != nil {
        panic(err)
    }

    // write to file for testing
    //file, err := os.Create("analyzes_hierarchy.txt")
    //if err != nil {
    //    fmt.Println("Analyzes hierarchy.", err)
    //    return
    //}
    //defer file.Close()

    // write to DB
    for a := range analyse_chan {

        // Insert analyse_type
        var type_id int64
        err = db.QueryRow("SELECT id FROM analyse_type WHERE type_name =? LIMIT 1;", a.Type).Scan(&type_id)
        if err == sql.ErrNoRows {
            res, _ := db.Exec("INSERT INTO analyse_type(type_name) values(?);", a.Type)
            type_id, err = res.LastInsertId()
        } else if err != nil {
            panic(err)
        }

        if a.Subtype != "" {
            // Insert analyse_subtype
            var subtype_id int64
            err = db.QueryRow("SELECT id FROM analyse_subtype WHERE subtype_name =? LIMIT 1;", a.Subtype).Scan(&subtype_id)
            if err == sql.ErrNoRows {
                res, _ := db.Exec("INSERT INTO analyse_subtype(subtype_name, type) values(?, ?);", a.Subtype, type_id)
                subtype_id, err = res.LastInsertId()
            } else if err != nil {
                panic(err)
            }
            // Insert analyse
            sql_insert := "INSERT INTO analyse(type, subtype, analyse_kind, decribe, prepare, indications, interpritation) values(?, ?, ?, ?, ?, ?, ?);"
            //sql_insert := "INSERT INTO analyse(type, subtype, analyse_kind) values(?, ?, ?);"
            res, _ := db.Exec(sql_insert, type_id, subtype_id, a.Kind, a.DetailDesc, a.DetailPrep, a.DetailIndic, a.DetailInterp)
            //res, _ := db.Exec(sql_insert, type_id, subtype_id, a.Kind)
            _, err := res.LastInsertId()
            if err != nil {
                panic(err)
            }
        } else {
            // Insert analyse
            sql_insert := "INSERT INTO analyse(type, analyse_kind, decribe, prepare, indications, interpritation) values(?, ?, ?, ?, ?, ?);"
            //sql_insert := "INSERT INTO analyse(type, subtype, analyse_kind) values(?, ?, ?);"
            res, _ := db.Exec(sql_insert, type_id, a.Kind, a.DetailDesc, a.DetailPrep, a.DetailIndic, a.DetailInterp)
            //res, _ := db.Exec(sql_insert, type_id, subtype_id, a.Kind)
            _, err := res.LastInsertId()
            if err != nil {
                panic(err)
            }
        }

        // write to file for testing.
        //file.WriteString("--" + a.Type + "\n")
        //file.WriteString("----" + a.Subtype + "\n")
        //file.WriteString("------" + a.Kind + "\n")
        //file.WriteString("--------" + a.DetailDesc + "\n")
        //file.WriteString("--------" + a.DetailPrep + "\n")
        //file.WriteString("--------" + a.DetailIndic + "\n")
        //file.WriteString("--------" + a.DetailInterp + "\n")
    }
}


func getAnalyseDetailByURL(analyse Analyse, analyse_chan chan Analyse) {
    defer parser_wg.Done()

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
    //details_html_rd, _ := os.Open("test_html/test_156_28932.html")
    //defer details_html_rd.Close()

    details_html_tree, _ :=	html.Parse(details_html_rd)
    details_node := getDetailsNode(details_html_tree)
    if details_node == nil {return}
    details_subnodes_raw := details_node.FirstChild.Data

    details_subnodes_htmls := getDetailsSubnodesHTMLs(details_subnodes_raw)

    for det, d_html := range details_subnodes_htmls {
        detail, _ := html.Parse(strings.NewReader(d_html))
        detail_text := strings.Join(getDetailsNodeText(nil, detail), "")
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
