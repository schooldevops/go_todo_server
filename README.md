# TODO Project 

## Install Package 

### Gorilla 

[Gorilla](https://github.com/gorilla/mux)

```
go get "github.com/gorilla/mux"
```

### MariaDB 

```
go get "github.com/go-sql-driver/mysql"
```

## DB 접속정보 

```
"boarduser:boarduser123@tcp(127.0.0.1:3306)/simpleboard?charset=utf8&parseTime=True&loc=Local"
```


## Reference

[https://gowebexamples.com/routes-using-gorilla-mux/](https://gowebexamples.com/routes-using-gorilla-mux/)

[https://techinscribed.com/different-approaches-to-pass-database-connection-into-controllers-in-golang/](https://techinscribed.com/different-approaches-to-pass-database-connection-into-controllers-in-golang/)

[https://mariadb.com/ko/resources/blog/using-go-with-mariadb/](https://mariadb.com/ko/resources/blog/using-go-with-mariadb/)

[https://forum.golangbridge.org/t/best-practice-to-use-go-sql-driver-mysql-package/8028](https://forum.golangbridge.org/t/best-practice-to-use-go-sql-driver-mysql-package/8028)




## Docker Push

### Docker Build 
docker build -t XXXXXXX.dkr.ecr.ap-northeast-2.amazonaws.com/todoapp:1.0 .

### Create ECR
aws ecr create-repository \
    --repository-name todoapp \
    --image-scanning-configuration scanOnPush=true \
    --region ap-northeast-2

### ECR Login
aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin XXXXXXX.dkr.ecr.ap-northeast-2.amazonaws.com

### ECS Push 
docker push XXXXXXX.dkr.ecr.ap-northeast-2.amazonaws.com/todoapp:1.0

### Delete Image 
aws ecr batch-delete-image \
      --repository-name todoapp \
      --image-ids imageTag=latest

### Delete Repository
aws ecr delete-repository \
      --repository-name todoapp \
      --force      