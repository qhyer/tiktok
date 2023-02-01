cd ./api
go build main.go
cd ..

cd comment
sh build.sh
cd ..

cd ./favorite
sh build.sh
cd ..

cd ./feed
sh build.sh
cd ..

cd ./publish
sh build.sh
cd ..

cd ./relation
sh build.sh
cd ..

cd ./user
sh build.sh
cd ..

cd ./message
sh build.sh