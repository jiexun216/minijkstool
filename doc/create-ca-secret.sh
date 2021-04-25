#!/bin/bash -x

kubectl delete secret ca-demo

rm -vf ca.crt truststore.jks

kubectl -n cdp-dce get secret default-token-swf8x -o yaml \
	| yq read - data[ca\.crt] \
	| base64 -d > ca.crt

# 将第三方给的cer文件加入自己的秘钥库
 keytool -import  -storepass changeit -noprompt \
	-alias cert0 -file ca.crt \
	-storetype JKS -keystore truststore.jks

docker run --rm -v "$(pwd)":/drop openjdk:slim-buster \
       keytool -import  -storepass changeit -noprompt \
       -alias cert0 -file /drop/ca.crt \
       -storetype JKS -keystore /drop/truststore.jks

kubectl create secret generic ca-demo --from-file=cacerts=truststore.jks

kubectl get secret ca-demo -o yaml
