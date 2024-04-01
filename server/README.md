
Installation
1. go get github.com/rakyll/statik
2. go get -u github.com/go-swagger/go-swagger/cmd/swagger #For downloading the Go Swagger to create the spec using the swagger comments.


Steps to follow
1. Add the Swagger Comments to the API added or updated using documention at https://goswagger.io/use/spec.html.
2. Run GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models  from the root directory.
3. cd metis-seq/themis/server
4. Replace the Swagger.yaml file inside swagger-ui directory with the swagger.yaml newly generated in root directly in step 2
5. cd metis-seq/themis/server && statik -src=./swagger-ui
6. cd metis-seq/themis && make build
7. cd metis-seq/themis && make run-server

Steps to follow for updated swagger-ui without using go-swagger
1. Add the Swagger Comments to the API added or updated using documention at https://goswagger.io/use/spec.html.
2. Run GO111MODULE=off swagger generate spec -o ./swagger.yaml --scan-models  from the root directory.
3. Copy zip file of source code from https://github.com/swagger-api/swagger-ui/releases
4. Unzip the zip file. Copy the contents of dist/ from the zip to themis/server/swagger-ui/
5. Convert the themis/server/swagger-ui/swagger.yaml to JSON format and place it in the same directory as the swagger.yaml file.
6. In themis/server/swagger-ui/swagger-initializer.js change `url: "./swagger.json"`,
7. cd metis-seq/themis/server && statik -src=./swagger-ui

Visit http://localhost:1317/swagger-ui/ 


Reference
- https://github.com/rakyll/statik
