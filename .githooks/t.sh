#!/bin/bash
Precent=75.0


#检查质量
function checkQ(){
   for FILE in `git diff --name-only --cached`; do
	if [[ $FILE == *.go ]]; then
	   #echo "golangci-lint run --disable=typecheck $FILE -c ./.githooks/.golangci.yml"
	   golangci-lint run --disable=typecheck $FILE -c ./.githooks/.golangci.yml
	   if [[ $? != 0 ]]; then
              exit 1
           fi
	fi
   done 
}

#检查特定单词
function checkWord(){
for FILE in `git diff --diff-filter=ACMTX --name-only --cached`; do
    if [[ $FILE == *.go ]]; then
    	grep 'TODO:\|debugger\|console.log\|alert(' $FILE 2>&1 >/dev/null
    	if [ $? -eq 0 ]; then
        	echo $FILE '包含，TODO: or debugger or console.log，请删除后再提交'
        	exit 1
    	fi
    fi
done
}

#单测检查
function checkTest(){
  res=`go test $@ -coverprofile=cover.out`
  IFS=$'\n\n'
  for line in $res;do
    	  str=`echo $line | awk '{print $5}'`
          if [[ $str = *%* ]]; then #判断是否有百分比数值
	      p=`echo ${str%*\%}`
	      isPass=`echo "$p >= $Precent" | bc`
	      if [[ $isPass -ne 1 ]]; then 
			echo -e "代码覆盖率不符合$Precent% 标准\n$line"
			exit 1
              fi
	  fi
  done
}

#checkTest "gomybatis/util" "gomybatis/cmd"
#checkWord
#checkQ
