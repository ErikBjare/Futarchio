#!/bin/bash

re_js="site/src/.*[^(\.min)]\.js$"
re_html="site/src/[^/]*\.html$"
re_go="[^(_test)]\.go$"
re_go_tests="_test\.go$"
re_md="\.md$"

echo -n "Lines of HTML: "
cat `find src | grep $re_html` | wc -l

echo -n "Lines of JS: "
cat `find src | grep $re_js` | wc -l

echo -n "Lines of Go (excluding tests): "
cat `find src | grep $re_go` | wc -l
echo -n "Lines of Go tests: "
cat `find src | grep $re_go_tests` | wc -l

echo -n "Lines of Markdown: "
cat `find wiki | grep $re_md` | wc -l

#echo ""
#echo "Verbose:"
#echo "JS: $(echo `find | grep $re_js`)"
#echo "HTML: $(echo `find | grep $re_html`)"
#echo "Go: $(echo `find | grep $re_go`)"
#echo "Go: $(echo `find | grep $re_go_tests`)"
