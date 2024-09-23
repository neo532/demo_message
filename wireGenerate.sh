#/bin/bash

fileName='wireProviderSet.go'

for d in $(find . -name $fileName);
do
    cntOri=$(cat $d)
    dir=$(dirname $d)

    package=$(basename $dir)
    news=$(grep -r '^func New' $dir|sort|awk '{print $2}'|awk -F'(' '{printf "\\t%s,\\n",$1}')
    tpl=$(echo -e "// generate by wireGenerate.sh with '^func New' in on package\npackage $package\n\nimport (\n\t\"github.com/google/wire\"\n)\n\nvar ProviderSet = wire.NewSet(\n$news)")

    tplC=$(echo $tpl|tr -d '\n'|sed -E 's/[[:space:]]*//g')
    cntOriC=$(echo $cntOri|tr -d '\n'|sed -E 's/[[:space:]]*//g')
      # echo "tplC:"$tplC
      # echo "oriC:"$cntOriC
    if [ $tplC"x" != $cntOriC"x" ] || [ $(wc -l $d|awk '{print $1}') == "1" ]; then
        echo -e $d
        echo -e "$tpl"|sed "s/-e//g" > $dir/$fileName
    fi
    #echo $(stat  $d)
done
