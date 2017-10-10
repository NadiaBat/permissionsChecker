package main

import (
	"database/sql"

	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type MySQLConnectionConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	Base       string `yaml:"base"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Parameters struct {
		MaxIdleConns int `yaml:"max_idle_conns"`
		MaxOpenConns int `yaml:"max_open_conns"`
	} `yaml:"parameters"`
}

// @TODO 2
var rulesDictionary = map[string]string{
	"a:2:{s:9:\"paramsKey\";s:12:\"isCommercial\";s:4:\"Data\";a:1:{i:0;i:0;}}":     "{\"paramsKey\":\"isCommercial\",\"Data\":[\"0\"]}",
	"a:2:{s:9:\"paramsKey\";s:12:\"isCommercial\";s:4:\"Data\";a:1:{i:0;i:1;}}":     "{\"paramsKey\":\"isCommercial\",\"Data\":[\"1\"]}",
	"a:2:{s:9:\"paramsKey\";s:3:\"pid\";s:4:\"Data\";a:1:{i:0;s:6:\"369550\";}}":    "{\"paramsKey\":\"pid\",\"Data\":[\"369550\"]}",
	"a:2:{s:9:\"paramsKey\";s:3:\"pid\";s:4:\"Data\";a:1:{i:0;s:8:\"14338667\";}}":  "{\"paramsKey\":\"pid\",\"Data\":[\"14338667\"]}",
	"a:2:{s:9:\"paramsKey\";s:3:\"pid\";s:4:\"Data\";a:1:{i:0;s:8:\"14338727\";}}":  "{\"paramsKey\":\"pid\",\"Data\":[\"14338727\"]}",
	"a:2:{s:9:\"paramsKey\";s:3:\"pid\";s:4:\"Data\";a:1:{i:0;s:8:\"14338747\";}}":  "{\"paramsKey\":\"pid\",\"Data\":[\"14338747\"]}",
	"a:2:{s:9:\"paramsKey\";s:3:\"pid\";s:4:\"Data\";a:1:{i:0;s:9:\"145919821\";}}": "{\"paramsKey\":\"pid\",\"Data\":[\"145919821\"]}",
	"a:2:{s:9:\"paramsKey\";s:3:\"pid\";s:4:\"Data\";a:1:{i:0;s:9:\"152645602\";}}": "{\"paramsKey\":\"pid\",\"Data\":[\"152645602\"]}",
	"a:2:{s:9:\"paramsKey\";s:3:\"pid\";s:4:\"Data\";a:1:{i:0;s:9:\"200132743\";}}": "{\"paramsKey\":\"pid\",\"Data\":[\"200132743\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";N;}":                     "{\"paramsKey\":\"paramsKey\",\"Data\":null}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:0:{}}":                 "{\"paramsKey\":\"paramsKey\",\"Data\":[]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:10:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:15:\"flags[yarabota]\";i:3;s:4:\"type\";i:4;s:9:\"mainPhoto\";i:5;s:6:\"images\";i:6;s:10:\"sourceName\";i:7;s:5:\"theme\";i:8;s:7:\"authors\";i:9;s:17:\"isCommentsAllowed\";}}":                                                                                                                                                                                                       "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"flags[yarabota]\",\"type\",\"mainPhoto\",\"images\",\"sourceName\",\"theme\",\"authors\",\"isCommentsAllowed\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:10:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:4:\"type\";i:3;s:6:\"images\";i:4;s:10:\"sourceName\";i:5;s:5:\"theme\";i:6;s:5:\"links\";i:7;s:7:\"authors\";i:8;s:17:\"isCommentsAllowed\";i:9;s:4:\"tags\";}}":                                                                                                                                                                                                                       "{\"paramsKey\":\"paramsKey\",\"Data\":[\"lead\",\"subheader\",\"type\",\"images\",\"sourceName\",\"theme\",\"links\",\"authors\",\"isCommentsAllowed\",\"tags\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:11:{i:0;s:13:\"flags[tomain]\";i:1;s:12:\"flags[torss]\";i:2;s:5:\"isBaa\";i:3;s:9:\"mainPhoto\";i:4;s:10:\"sourceName\";i:5;s:8:\"category\";i:6;s:5:\"theme\";i:7;s:5:\"links\";i:8;s:7:\"authors\";i:9;s:4:\"tags\";i:10;s:9:\"copyright\";}}":                                                                                                                                                                                        "{\"paramsKey\":\"paramsKey\",\"Data\":\"flags[tomain]\",\"flags[torss]\",\"isBaa\",\"mainPhoto\",\"sourceName\",\"category\",\"theme\",\"links\",\"authors\",\"tags\",\"copyright\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:12:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:12:\"flags[torss]\";i:3;s:4:\"type\";i:4;s:6:\"images\";i:5;s:10:\"sourceName\";i:6;s:5:\"theme\";i:7;s:5:\"links\";i:8;s:7:\"authors\";i:9;s:17:\"isCommentsAllowed\";i:10;s:4:\"tags\";i:11;s:8:\"crmOrder\";}}":                                                                                                                                                                      "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"flags[torss]\",\"type\",\"images\",\"sourceName\",\"theme\",\"links\",\"authors\",\"isCommentsAllowed\",\"tags\",\"crmOrder\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:12:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:5:\"isBaa\";i:3;s:4:\"type\";i:4;s:9:\"mainPhoto\";i:5;s:21:\"socialBackgroundImage\";i:6;s:6:\"images\";i:7;s:10:\"sourceName\";i:8;s:5:\"theme\";i:9;s:7:\"authors\";i:10;s:17:\"isCommentsAllowed\";i:11;s:8:\"crmOrder\";}}":                                                                                                                                                        "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"isBaa\",\"type\",\"mainPhoto\",\"socialBackgroundImage\",\"images\",\"sourceName\",\"theme\",\"authors\",\"isCommentsAllowed\",\"crmOrder\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:12:{i:0;s:9:\"subheader\";i:1;s:12:\"imagesAuthor\";i:2;s:13:\"flags[tomain]\";i:3;s:9:\"mainPhoto\";i:4;s:6:\"images\";i:5;s:6:\"videos\";i:6;s:8:\"category\";i:7;s:8:\"keywords\";i:8;s:5:\"polls\";i:9;s:5:\"links\";i:10;s:7:\"authors\";i:11;s:17:\"isCommentsAllowed\";}}":                                                                                                                                                        "{\"paramsKey\":\"paramsKey\",\"Data\":\"subheader\",\"imagesAuthor\",\"flags[tomain]\",\"mainPhoto\",\"images\",\"videos\",\"category\",\"keywords\",\"polls\",\"links\",\"authors\",\"isCommentsAllowed\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:12:{i:0;s:9:\"subheader\";i:1;s:12:\"imagesAuthor\";i:2;s:13:\"flags[tomain]\";i:3;s:9:\"mainPhoto\";i:4;s:6:\"images\";i:5;s:8:\"category\";i:6;s:8:\"keywords\";i:7;s:5:\"polls\";i:8;s:5:\"links\";i:9;s:7:\"authors\";i:10;s:17:\"isCommentsAllowed\";i:11;s:8:\"crmOrder\";}}":                                                                                                                                                      "{\"paramsKey\":\"paramsKey\",\"Data\":\"subheader\",\"imagesAuthor\",\"flags[tomain]\",\"mainPhoto\",\"images\",\"category\",\"keywords\",\"polls\",\"links\",\"authors\",\"isCommentsAllowed\",\"crmOrder\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:13:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:10:\"htmlButton\";i:3;s:5:\"feeds\";i:4;s:4:\"type\";i:5;s:6:\"images\";i:6;s:6:\"videos\";i:7;s:10:\"sourceName\";i:8;s:5:\"theme\";i:9;s:5:\"links\";i:10;s:7:\"authors\";i:11;s:17:\"isCommentsAllowed\";i:12;s:4:\"tags\";}}":                                                                                                                                                       "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"htmlButton\",\"feeds\",\"type\",\"images\",\"videos\",\"sourceName\",\"theme\",\"links\",\"authors\",\"isCommentsAllowed\",\"tags\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:14:{i:0;s:16:\"32.autoBlock.fri\";i:1;s:16:\"32.autoBlock.mon\";i:2;s:16:\"32.autoBlock.sat\";i:3;s:16:\"32.autoBlock.sun\";i:4;s:16:\"32.autoBlock.thu\";i:5;s:16:\"32.autoBlock.tue\";i:6;s:16:\"32.autoBlock.wed\";i:7;s:15:\"32.domBlock.fri\";i:8;s:15:\"32.domBlock.mon\";i:9;s:15:\"32.domBlock.sat\";i:10;s:15:\"32.domBlock.sun\";i:11;s:15:\"32.domBlock.thu\";i:12;s:15:\"32.domBlock.tue\";i:13;s:15:\"32.domBlock.wed\";}}": "{\"paramsKey\":\"paramsKey\",\"Data\":\"32.autoBlock.fri\",\"32.autoBlock.mon\",\"32.autoBlock.sat\",\"32.autoBlock.sun\",\"32.autoBlock.thu\",\"32.autoBlock.tue\",\"32.autoBlock.wed\",\"32.domBlock.fri\",\"32.domBlock.mon\",\"32.domBlock.sat\",\"32.domBlock.sun\",\"32.domBlock.thu\",\"32.domBlock.tue\",\"32.domBlock.wed\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:14:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:12:\"flags[torss]\";i:3;s:4:\"type\";i:4;s:6:\"images\";i:5;s:6:\"videos\";i:6;s:10:\"sourceName\";i:7;s:8:\"category\";i:8;s:5:\"theme\";i:9;s:5:\"links\";i:10;s:7:\"authors\";i:11;s:17:\"isCommentsAllowed\";i:12;s:4:\"tags\";i:13;s:8:\"crmOrder\";}}":                                                                                                                            "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"flags[torss]\",\"type\",\"images\",\"videos\",\"sourceName\",\"category\",\"theme\",\"links\",\"authors\",\"isCommentsAllowed\",\"tags\",\"crmOrder\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:14:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:5:\"feeds\";i:3;s:13:\"flags[tomain]\";i:4;s:5:\"isBaa\";i:5;s:4:\"type\";i:6;s:9:\"mainPhoto\";i:7;s:21:\"socialBackgroundImage\";i:8;s:6:\"author\";i:9;s:8:\"category\";i:10;s:5:\"theme\";i:11;s:7:\"authors\";i:12;s:17:\"isCommentsAllowed\";i:13;s:9:\"customUrl\";}}":                                                                                                           "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"feeds\",\"flags[tomain]\",\"isBaa\",\"type\",\"mainPhoto\",\"socialBackgroundImage\",\"author\",\"category\",\"theme\",\"authors\",\"isCommentsAllowed\",\"customUrl\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:15:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:10:\"htmlButton\";i:3;s:13:\"flags[tomain]\";i:4;s:12:\"flags[totop]\";i:5;s:12:\"flags[torss]\";i:6;s:4:\"type\";i:7;s:9:\"mainPhoto\";i:8;s:6:\"images\";i:9;s:8:\"category\";i:10;s:5:\"theme\";i:11;s:5:\"polls\";i:12;s:7:\"authors\";i:13;s:17:\"isCommentsAllowed\";i:14;s:15:\"commercialLabel\";}}":                                                                            "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"htmlButton\",\"flags[tomain]\",\"flags[totop]\",\"flags[torss]\",\"type\",\"mainPhoto\",\"images\",\"category\",\"theme\",\"polls\",\"authors\",\"isCommentsAllowed\",\"commercialLabel\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:15:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:13:\"flags[tomain]\";i:3;s:12:\"flags[totop]\";i:4;s:12:\"flags[torss]\";i:5;s:4:\"type\";i:6;s:9:\"mainPhoto\";i:7;s:6:\"images\";i:8;s:10:\"sourceName\";i:9;s:5:\"theme\";i:10;s:5:\"links\";i:11;s:7:\"authors\";i:12;s:17:\"isCommentsAllowed\";i:13;s:4:\"tags\";i:14;s:9:\"copyright\";}}":                                                                                       "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"flags[tomain]\",\"flags[totop]\",\"flags[torss]\",\"type\",\"mainPhoto\",\"images\",\"sourceName\",\"theme\",\"links\",\"authors\",\"isCommentsAllowed\",\"tags\",\"copyright\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:16:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:10:\"htmlButton\";i:3;s:13:\"flags[tomain]\";i:4;s:15:\"flags[yarabota]\";i:5;s:5:\"isBaa\";i:6;s:4:\"type\";i:7;s:9:\"mainPhoto\";i:8;s:6:\"images\";i:9;s:6:\"videos\";i:10;s:10:\"sourceName\";i:11;s:5:\"theme\";i:12;s:7:\"authors\";i:13;s:17:\"isCommentsAllowed\";i:14;s:9:\"customUrl\";i:15;s:9:\"copyright\";}}":                                                             "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"htmlButton\",\"flags[tomain]\",\"flags[yarabota]\",\"isBaa\",\"type\",\"mainPhoto\",\"images\",\"videos\",\"sourceName\",\"theme\",\"authors\",\"isCommentsAllowed\",\"customUrl\",\"copyright\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:17:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:10:\"htmlButton\";i:3;s:13:\"flags[tomain]\";i:4;s:5:\"isBaa\";i:5;s:4:\"type\";i:6;s:9:\"mainPhoto\";i:7;s:6:\"images\";i:8;s:6:\"videos\";i:9;s:10:\"sourceName\";i:10;s:5:\"theme\";i:11;s:7:\"authors\";i:12;s:17:\"isCommentsAllowed\";i:13;s:15:\"commercialLabel\";i:14;s:9:\"customUrl\";i:15;s:9:\"copyright\";i:16;s:8:\"crmOrder\";}}":                                       "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"htmlButton\",\"flags[tomain]\",\"isBaa\",\"type\",\"mainPhoto\",\"images\",\"videos\",\"sourceName\",\"theme\",\"authors\",\"isCommentsAllowed\",\"commercialLabel\",\"customUrl\",\"copyright\",\"crmOrder\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:17:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:13:\"flags[tomain]\";i:3;s:12:\"flags[torss]\";i:4;s:13:\"flags[n1spam]\";i:5;s:4:\"type\";i:6;s:9:\"mainPhoto\";i:7;s:6:\"images\";i:8;s:10:\"sourceName\";i:9;s:5:\"theme\";i:10;s:14:\"headerKeywords\";i:11;s:5:\"links\";i:12;s:7:\"authors\";i:13;s:17:\"isCommentsAllowed\";i:14;s:4:\"tags\";i:15;s:9:\"customUrl\";i:16;s:9:\"copyright\";}}":                                  "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"flags[tomain]\",\"flags[torss]\",\"flags[n1spam]\",\"type\",\"mainPhoto\",\"images\",\"sourceName\",\"theme\",\"headerKeywords\",\"links\",\"authors\",\"isCommentsAllowed\",\"tags\",\"customUrl\",\"copyright\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:17:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:13:\"flags[tomain]\";i:3;s:12:\"flags[tovip]\";i:4;s:4:\"type\";i:5;s:9:\"mainPhoto\";i:6;s:9:\"mainVideo\";i:7;s:6:\"images\";i:8;s:6:\"videos\";i:9;s:10:\"sourceName\";i:10;s:8:\"category\";i:11;s:5:\"theme\";i:12;s:5:\"polls\";i:13;s:17:\"isCommentsAllowed\";i:14;s:15:\"commercialLabel\";i:15;s:9:\"customUrl\";i:16;s:9:\"copyright\";}}":                                   "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"flags[tomain]\",\"flags[tovip]\",\"isBaa\",\"type\",\"mainPhoto\",\"mainVideo\",\"images\",\"sourceName\",\"author\",\"category\",\"theme\",\"polls\",\"commercialLabel\",\"customUrl\",\"copyright\",\"crmOrder\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:18:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:13:\"flags[tomain]\";i:3;s:12:\"flags[tovip]\";i:4;s:5:\"isBaa\";i:5;s:4:\"type\";i:6;s:9:\"mainPhoto\";i:7;s:9:\"mainVideo\";i:8;s:6:\"images\";i:9;s:10:\"sourceName\";i:10;s:6:\"author\";i:11;s:8:\"category\";i:12;s:5:\"theme\";i:13;s:5:\"polls\";i:14;s:15:\"commercialLabel\";i:15;s:9:\"customUrl\";i:16;s:9:\"copyright\";i:17;s:8:\"crmOrder\";}}":                          "{\"paramsKey\":\"paramsKey\",\"Data\":[\"flags[n1spam]\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:1:{i:0;s:13:\"flags[n1spam]\";}}":    "{\"paramsKey\":\"paramsKey\",\"Data\":[\"flags[comPlace7]\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:1:{i:0;s:16:\"flags[comPlace7]\";}}": "{\"paramsKey\":\"paramsKey\",\"Data\":[\"lead\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:1:{i:0;s:4:\"lead\";}}":              "{\"paramsKey\":\"paramsKey\",\"Data\":[\"authors\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:1:{i:0;s:7:\"authors\";}}":           "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"imagesAuthor\",\"flags[tomain]\",\"flags[fixonmain]\",\"flags[toslider]\",\"flags[photorep]\",\"isBaa\",\"type\",\"mainPhoto\",\"images\",\"sourceName\",\"category\",\"theme\",\"links\",\"authors\",\"isCommentsAllowed\",\"commercialLabel\",\"customUrl\",\"copyright\",\"crmOrder\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:21:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:12:\"imagesAuthor\";i:3;s:13:\"flags[tomain]\";i:4;s:16:\"flags[fixonmain]\";i:5;s:15:\"flags[toslider]\";i:6;s:15:\"flags[photorep]\";i:7;s:5:\"isBaa\";i:8;s:4:\"type\";i:9;s:9:\"mainPhoto\";i:10;s:6:\"images\";i:11;s:10:\"sourceName\";i:12;s:8:\"category\";i:13;s:5:\"theme\";i:14;s:5:\"links\";i:15;s:7:\"authors\";i:16;s:17:\"isCommentsAllowed\";i:17;s:15:\"commercialLabel\";i:18;s:9:\"customUrl\";i:19;s:9:\"copyright\";i:20;s:8:\"crmOrder\";}}":                                                                                                                                                                            "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"imagesAuthor\",\"flags[tomain]\",\"flags[totop]\",\"flags[torss]\",\"flags[toslider]\",\"flags[photorep]\",\"type\",\"mainPhoto\",\"socialBackgroundImage\",\"socialBackgroundImageHeader\",\"images\",\"sourceName\",\"category\",\"theme\",\"headerKeywords\",\"links\",\"authors\",\"isCommentsAllowed\",\"customUrl\",\"copyright\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:22:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:12:\"imagesAuthor\";i:3;s:13:\"flags[tomain]\";i:4;s:12:\"flags[totop]\";i:5;s:12:\"flags[torss]\";i:6;s:15:\"flags[toslider]\";i:7;s:15:\"flags[photorep]\";i:8;s:4:\"type\";i:9;s:9:\"mainPhoto\";i:10;s:21:\"socialBackgroundImage\";i:11;s:27:\"socialBackgroundImageHeader\";i:12;s:6:\"images\";i:13;s:10:\"sourceName\";i:14;s:8:\"category\";i:15;s:5:\"theme\";i:16;s:14:\"headerKeywords\";i:17;s:5:\"links\";i:18;s:7:\"authors\";i:19;s:17:\"isCommentsAllowed\";i:20;s:9:\"customUrl\";i:21;s:9:\"copyright\";}}":                                                                                                                 "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"flags[tomain]\",\"flags[totop]\",\"flags[torss]\",\"type\",\"mainPhoto\",\"mainVideo\",\"images\",\"videos\",\"sourceName\",\"category\",\"theme\",\"headerKeywords\",\"links\",\"authors\",\"goodNewsText\",\"isCommentsAllowed\",\"tags\",\"commercialLabel\",\"customUrl\",\"copyright\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:22:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:13:\"flags[tomain]\";i:3;s:12:\"flags[totop]\";i:4;s:12:\"flags[torss]\";i:5;s:4:\"type\";i:6;s:9:\"mainPhoto\";i:7;s:9:\"mainVideo\";i:8;s:6:\"images\";i:9;s:6:\"videos\";i:10;s:10:\"sourceName\";i:11;s:8:\"category\";i:12;s:5:\"theme\";i:13;s:14:\"headerKeywords\";i:14;s:5:\"links\";i:15;s:7:\"authors\";i:16;s:12:\"goodNewsText\";i:17;s:17:\"isCommentsAllowed\";i:18;s:4:\"tags\";i:19;s:15:\"commercialLabel\";i:20;s:9:\"customUrl\";i:21;s:9:\"copyright\";}}":                                                                                                                                                                "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"htmlButton\",\"feeds\",\"flags[tomain]\",\"flags[totop]\",\"flags[torss]\",\"flags[tovip]\",\"type\",\"mainPhoto\",\"mainVideo\",\"images\",\"videos\",\"sourceName\",\"category\",\"theme\",\"headerKeywords\",\"links\",\"authors\",\"goodNewsText\",\"isCommentsAllowed\",\"tags\",\"commercialLabel\",\"copyright\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:24:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:10:\"htmlButton\";i:3;s:5:\"feeds\";i:4;s:13:\"flags[tomain]\";i:5;s:12:\"flags[totop]\";i:6;s:12:\"flags[torss]\";i:7;s:12:\"flags[tovip]\";i:8;s:4:\"type\";i:9;s:9:\"mainPhoto\";i:10;s:9:\"mainVideo\";i:11;s:6:\"images\";i:12;s:6:\"videos\";i:13;s:10:\"sourceName\";i:14;s:8:\"category\";i:15;s:5:\"theme\";i:16;s:14:\"headerKeywords\";i:17;s:5:\"links\";i:18;s:7:\"authors\";i:19;s:12:\"goodNewsText\";i:20;s:17:\"isCommentsAllowed\";i:21;s:4:\"tags\";i:22;s:15:\"commercialLabel\";i:23;s:9:\"copyright\";}}":                                                                                                                "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"htmlButton\",\"flags[tomain]\",\"flags[totop]\",\"flags[torss]\",\"flags[tovip]\",\"type\",\"mainPhoto\",\"socialBackgroundImage\",\"socialBackgroundImageHeader\",\"images\",\"videos\",\"sourceName\",\"category\",\"theme\",\"headerKeywords\",\"links\",\"authors\",\"goodNewsText\",\"isCommentsAllowed\",\"tags\",\"commercialLabel\",\"customUrl\",\"copyright\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:25:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:10:\"htmlButton\";i:3;s:13:\"flags[tomain]\";i:4;s:12:\"flags[totop]\";i:5;s:12:\"flags[torss]\";i:6;s:12:\"flags[tovip]\";i:7;s:4:\"type\";i:8;s:9:\"mainPhoto\";i:9;s:21:\"socialBackgroundImage\";i:10;s:27:\"socialBackgroundImageHeader\";i:11;s:6:\"images\";i:12;s:6:\"videos\";i:13;s:10:\"sourceName\";i:14;s:8:\"category\";i:15;s:5:\"theme\";i:16;s:14:\"headerKeywords\";i:17;s:5:\"links\";i:18;s:7:\"authors\";i:19;s:12:\"goodNewsText\";i:20;s:17:\"isCommentsAllowed\";i:21;s:4:\"tags\";i:22;s:15:\"commercialLabel\";i:23;s:9:\"customUrl\";i:24;s:9:\"copyright\";}}":                                                     "{\"paramsKey\":\"paramsKey\",\"Data\":\"lead\",\"subheader\",\"flags[tomain]\",\"flags[totop]\",\"flags[tofixedtop]\",\"flags[torss]\",\"flags[goodNews]\",\"type\",\"mainPhoto\",\"socialBackgroundImage\",\"socialBackgroundImageHeader\",\"mainVideo\",\"images\",\"sourceName\",\"author\",\"category\",\"theme\",\"headerKeywords\",\"polls\",\"links\",\"authors\",\"goodNewsText\",\"isCommentsAllowed\",\"tags\",\"commercialLabel\",\"customUrl\",\"copyright\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:27:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:13:\"flags[tomain]\";i:3;s:12:\"flags[totop]\";i:4;s:17:\"flags[tofixedtop]\";i:5;s:12:\"flags[torss]\";i:6;s:15:\"flags[goodNews]\";i:7;s:4:\"type\";i:8;s:9:\"mainPhoto\";i:9;s:21:\"socialBackgroundImage\";i:10;s:27:\"socialBackgroundImageHeader\";i:11;s:9:\"mainVideo\";i:12;s:6:\"images\";i:13;s:10:\"sourceName\";i:14;s:6:\"author\";i:15;s:8:\"category\";i:16;s:5:\"theme\";i:17;s:14:\"headerKeywords\";i:18;s:5:\"polls\";i:19;s:5:\"links\";i:20;s:7:\"authors\";i:21;s:12:\"goodNewsText\";i:22;s:17:\"isCommentsAllowed\";i:23;s:4:\"tags\";i:24;s:15:\"commercialLabel\";i:25;s:9:\"customUrl\";i:26;s:9:\"copyright\";}}": "{\"paramsKey\":\"paramsKey\",\"Data\":[\"lead\",\"subheader\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:2:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";}}":                                                                                                      "{\"paramsKey\":\"paramsKey\",\"Data\":[\"category\",\"theme\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:2:{i:0;s:8:\"category\";i:1;s:5:\"theme\";}}":                                                                                                      "{\"paramsKey\":\"paramsKey\",\"Data\":[\"socialBackgroundImage\",\"socialBackgroundImageHeader\",\"socialAlign\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:3:{i:0;s:21:\"socialBackgroundImage\";i:1;s:27:\"socialBackgroundImageHeader\";i:2;s:11:\"socialAlign\";}}":                                        "{\"paramsKey\":\"paramsKey\",\"Data\":[\"lead\",\"type\",\"mainPhoto\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:3:{i:0;s:4:\"lead\";i:1;s:4:\"type\";i:2;s:9:\"mainPhoto\";}}":                                                                                     "{\"paramsKey\":\"paramsKey\",\"Data\":[\"subheader\",\"htmlButton\",\"isCommentsAllowed\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:3:{i:0;s:9:\"subheader\";i:1;s:10:\"htmlButton\";i:2;s:17:\"isCommentsAllowed\";}}":                                                                "{\"paramsKey\":\"paramsKey\",\"Data\":[\"socialBackgroundImage\",\"socialBackgroundImageHeader\",\"category\",\"theme\",\"socialAlign\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:5:{i:0;s:21:\"socialBackgroundImage\";i:1;s:27:\"socialBackgroundImageHeader\";i:2;s:8:\"category\";i:3;s:5:\"theme\";i:4;s:11:\"socialAlign\";}}": "{\"paramsKey\":\"paramsKey\",\"Data\":[\"feeds\",\"flags[tomain]\",\"flags[comPlace3]\",\"flags[comPlace7]\",\"isCommentsAllowed\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:5:{i:0;s:5:\"feeds\";i:1;s:13:\"flags[tomain]\";i:2;s:16:\"flags[comPlace3]\";i:3;s:16:\"flags[comPlace7]\";i:4;s:17:\"isCommentsAllowed\";}}":     "{\"paramsKey\":\"paramsKey\",\"Data\":[\"lead\",\"subheader\",\"sourceName\",\"headerKeywords\",\"authors\",\"copyright\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"paramsKey\";s:4:\"Data\";a:6:{i:0;s:4:\"lead\";i:1;s:9:\"subheader\";i:2;s:10:\"sourceName\";i:3;s:14:\"headerKeywords\";i:4;s:7:\"authors\";i:5;s:9:\"copyright\";}}":        "{\"paramsKey\":\"type\",\"Data\":[\"video_of_the_day\",\"photo_of_the_day\"]}",
	"a:2:{s:9:\"paramsKey\";s:4:\"type\";s:4:\"Data\";a:2:{i:0;s:16:\"video_of_the_day\";i:1;s:16:\"photo_of_the_day\";}}":                                                                                      "{\"paramsKey\":\"region\",\"Data\":[1077]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:1077;}}":                                                                                                                                     "{\"paramsKey\":\"region\",\"Data\":[114160]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:114160;}}":                                                                                                                                   "{\"paramsKey\":\"region\",\"Data\":[123]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:123;}}":                                                                                                                                      "{\"paramsKey\":\"region\",\"Data\":[124]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:124;}}":                                                                                                                                      "{\"paramsKey\":\"region\",\"Data\":[138]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:138;}}":                                                                                                                                      "{\"paramsKey\":\"region\",\"Data\":[142982]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:142982;}}":                                                                                                                                   "{\"paramsKey\":\"region\",\"Data\":[142]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:142;}}":                                                                                                                                      "{\"paramsKey\":\"region\",\"Data\":[14]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:14;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[154]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:154;}}":                                                                                                                                      "{\"paramsKey\":\"region\",\"Data\":[155]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:155;}}":                                                                                                                                      "{\"paramsKey\":\"region\",\"Data\":[166]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:166;}}":                                                                                                                                      "{\"paramsKey\":\"region\",\"Data\":[16]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:16;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[170]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:170;}}":                                                                                                                                      "{\"paramsKey\":\"region\",\"Data\":[181490]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:181490;}}":                                                                                                                                   "{\"paramsKey\":\"region\",\"Data\":[182028]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:182028;}}":                                                                                                                                   "{\"paramsKey\":\"region\",\"Data\":[18]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:18;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[21]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:21;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[22]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:22;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[23]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:23;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[24]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:24;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[26]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:26;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[27]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:27;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[29]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:29;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[2]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:2;}}":                                                                                                                                        "{\"paramsKey\":\"region\",\"Data\":[30]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:30;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[31]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:31;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[32]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:32;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[33]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:33;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[34]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:34;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[35]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:35;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[36]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:36;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[38]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:38;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[39]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:39;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[42]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:42;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[31]}43403]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:43403;}}":                                                                                                                                    "{\"paramsKey\":\"region\",\"Data\":[43]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:43;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[44]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:44;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[45]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:45;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[46]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:46;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[47]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:47;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[48]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:48;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[51]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:51;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[52]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:52;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[53]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:53;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[54]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:54;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[55]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:55;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[56]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:56;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[57]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:57;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[58]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:58;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[59]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:59;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[60]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:60;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[61]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:61;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[62]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:62;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[63]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:63;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[64]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:64;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[65]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:66;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[66]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:67;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[67]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:68;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[68]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:69;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[69]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:70;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[70]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:71;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[71]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:72;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[72]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:73;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[73]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:74;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[74]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:75;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[75]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:76;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[76]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:86;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[86]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:89;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[89]}",
	"a:2:{s:9:\"paramsKey\";s:6:\"region\";s:4:\"Data\";a:1:{i:0;i:93;}}":                                                                                                                                       "{\"paramsKey\":\"region\",\"Data\":[93]}",
	"a:2:{s:9:\"paramsKey\";s:7:\"project\";s:4:\"Data\";a:1:{i:0;i:10;}}":                                                                                                                                      "{\"paramsKey\":\"project\",\"Data\":[10]}",
	"a:2:{s:9:\"paramsKey\";s:7:\"project\";s:4:\"Data\";a:1:{i:0;i:11;}}":                                                                                                                                      "{\"paramsKey\":\"project\",\"Data\":[11]}",
	"a:2:{s:9:\"paramsKey\";s:7:\"project\";s:4:\"Data\";a:1:{i:0;i:1;}}":                                                                                                                                       "{\"paramsKey\":\"project\",\"Data\":[1]}",
	"a:2:{s:9:\"paramsKey\";s:7:\"project\";s:4:\"Data\";a:1:{i:0;i:28;}}":                                                                                                                                      "{\"paramsKey\":\"project\",\"Data\":[28]}",
	"a:2:{s:9:\"paramsKey\";s:7:\"project\";s:4:\"Data\";a:1:{i:0;i:2;}}":                                                                                                                                       "{\"paramsKey\":\"project\",\"Data\":[2]}",
	"a:2:{s:9:\"paramsKey\";s:7:\"project\";s:4:\"Data\";a:1:{i:0;i:3;}}":                                                                                                                                       "{\"paramsKey\":\"project\",\"Data\":[3]}",
	"a:2:{s:9:\"paramsKey\";s:7:\"project\";s:4:\"Data\";a:1:{i:0;i:48;}}":                                                                                                                                      "{\"paramsKey\":\"project\",\"Data\":[48]}",
	"a:2:{s:9:\"paramsKey\";s:7:\"project\";s:4:\"Data\";a:1:{i:0;i:4;}}":                                                                                                                                       "{\"paramsKey\":\"project\",\"Data\":[4]}",
	"a:2:{s:9:\"paramsKey\";s:7:\"project\";s:4:\"Data\";a:1:{i:0;i:6;}}":                                                                                                                                       "{\"paramsKey\":\"project\",\"Data\":[6]}",
	"a:2:{s:9:\"paramsKey\";s:7:\"project\";s:4:\"Data\";a:1:{i:0;i:7;}}":                                                                                                                                       "{\"paramsKey\":\"project\",\"Data\":[7]}",
	"a:2:{s:9:\"paramsKey\";s:7:\"project\";s:4:\"Data\";a:1:{i:0;i:8;}}":                                                                                                                                       "{\"paramsKey\":\"project\",\"Data\":[8]}",
	"a:2:{s:9:\"paramsKey\";s:7:\"project\";s:4:\"Data\";a:1:{i:0;i:9;}}":                                                                                                                                       "{\"paramsKey\":\"project\",\"Data\":[9]}",
	"a:2:{s:9:\"paramsKey\";s:8:\"template\";s:4:\"Data\";a:5:{i:0;i:1;i:1;i:2;i:2;i:3;i:3;i:4;i:4;i:5;}}":                                                                                                      "{\"paramsKey\":\"template\",\"Data\":[1,2,3,4,5]}",
	"a:2:{s:9:\"paramsKey\";s:8:\"template\";s:4:\"Data\";a:6:{i:0;i:1;i:1;i:2;i:2;i:3;i:3;i:4;i:4;i:5;i:5;i:6;}":                                                                                               "{\"paramsKey\":\"template\",\"Data\":[1,2,3,4,5,6]}",
}

func NewMySQL(config *MySQLConnectionConfig) (*sql.DB, error) {
	db, err := openDBConnection(config)
	if err != nil {
		return nil, errors.Wrap(err, "Fail create mysql client")
	}
	return db, nil
}

func openDBConnection(config *MySQLConnectionConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.GetDSN())
	if err != nil {
		return nil, errors.Wrapf(err, "can't open mysql connection"+
			" for dsn \"%s\"", config.GetDSN())
	}

	db.SetMaxIdleConns(config.Parameters.MaxIdleConns)
	db.SetMaxOpenConns(config.Parameters.MaxOpenConns)

	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "can't ping mysql "+
			"after open connection")
	}
	return db, nil
}

func (config MySQLConnectionConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Base,
	)
}

func RefreshCache() error {
	err := RefreshAssignments()
	if err != nil {
		return errors.Wrap(err, "Assignments refreshing failed")
	}

	err = RefreshPermissionItems()
	if err != nil {
		return errors.Wrap(err, "Permission items refreshing failed")
	}

	err = RefreshParents()
	if err != nil {
		return errors.Wrap(err, "Parents refreshing failed")
	}

	return nil
}

func RefreshAssignments() error {
	var err error
	Cache.assignments.Lock()
	Cache.assignments.data, err = getAssignmentsFromDb()
	Cache.assignments.Unlock()
	return err
}

func RefreshPermissionItems() error {
	var err error
	Cache.permissionItems.Lock()
	Cache.permissionItems.data, err = getPermissionItemsFromDb()
	Cache.permissionItems.Unlock()
	return err
}

func RefreshParents() error {
	var err error
	Cache.parents.Lock()
	Cache.parents.data, err = getParentsFromDb()
	Cache.parents.Unlock()
	return err
}

func getAssignmentsFromDb() (Assignments, error) {
	// как узнать пустая ли у тебя выборка?
	rows, err := mysql.Query(
		"SELECT IFNULL(`item_name`, ''), " +
			"IFNULL(`user_id`, 0), " +
			"IFNULL(`Data`, '') " +
			"FROM `auth_assignment`",
	)

	if err != nil {
		return nil, errors.Wrap(err, "Get all users assignments failed.")
	}

	result := Assignments{}
	var currentRule string
	var rule Rule

	for rows.Next() {
		var aRow AssignmentRow
		err = rows.Scan(
			&aRow.ItemName,
			&aRow.UserId,
			&currentRule,
		)

		if err != nil {
			errors.Wrapf(err, "Assignment row scanning error.")
		}

		if aRow.UserId == 0 || len(aRow.ItemName) == 0 {
			errors.Wrap(err, "Empty userId or itemName for assignment row.")
			continue
		}

		_, exist := result[aRow.UserId]
		if !exist {
			result[aRow.UserId] = UserAssignment{
				UserId: aRow.UserId,
				Items:  make(map[string]Assignment),
			}
		}

		rule, err = getRuleFromSerialized(currentRule)
		if err != nil {
			return nil, errors.Wrapf(err, "Unserialize currentRule failed. Rule was \"%s\"", currentRule)
		}

		result[aRow.UserId].Items[aRow.ItemName] = Assignment{
			ItemName: aRow.ItemName,
			Rule:     rule,
		}
	}

	return result, err
}

// @TODO 1
func getPermissionItemsFromDb() (PermissionItems, error) {
	res, err := mysql.Query(
		"SELECT IFNULL(`name`, ''), " +
			"IFNULL(`type`, 0), " +
			"IFNULL(`Data`, '') " +
			"FROM `auth_item`",
	)

	if err != nil {
		return nil, errors.Wrapf(err, "Auth items getting failed.")
	}

	currentName := ""
	currentType := 0
	currentRule := ""

	var rule Rule

	items := PermissionItems{}

	for res.Next() {
		var currentErr error
		currentErr = res.Scan(&currentName, &currentType, &currentRule)
		if currentErr != nil {
			err = errors.Wrap(err, "Auth item row scanning error.")
			continue
		}

		if len(currentName) == 0 {
			err = errors.Wrap(err, "Auth item paramsKey is empty.")
			continue
		}

		rule, currentErr = getRuleFromSerialized(currentRule)
		if currentErr != nil {
			err = errors.Wrapf(err, "Rule json decode error. Rule was \"%s\"", currentRule)
			continue
		}

		items[currentName] = PermissionItem{
			Name:     currentName,
			ItemType: currentType,
			Rule:     rule,
		}
	}

	return items, err
}

// @TODO 1
func getParentsFromDb() (AllParents, error) {
	res, err := mysql.Query("SELECT `child`, `parent` FROM `auth_item_child`")

	if err != nil {
		return nil, errors.Wrapf(err, "Parents getting failed.")
	}

	currentChild := ""
	currentParent := ""
	parents := AllParents{}
	for res.Next() {
		err := res.Scan(&currentChild, &currentParent)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"Auth item row scanning error with child %s and parent %s.",
				currentChild,
				currentParent,
			)
		}

		parents[currentChild] = append(parents[currentChild], currentParent)
	}

	return parents, nil
}

func getRuleFromSerialized(rule string) (Rule, error) {
	jsonRule, isExists := rulesDictionary[rule]
	fmt.Println(jsonRule)
	if !isExists {
		return Rule{}, nil
	}

	result := Rule{}
	err := json.Unmarshal([]byte(jsonRule), &result)

	return result, err
}
