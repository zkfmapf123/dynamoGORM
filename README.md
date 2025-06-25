# DynamoGORM

DynamoGORM은 Go 언어로 작성된 DynamoDB ORM 라이브러리입니다. AWS SDK v2를 기반으로 하며, 제네릭을 활용한 타입 안전한 데이터 조작을 제공합니다.

## 📦 설치

```bash
go get github.com/zkfmapf123/dynamoGORM
```

## 🔧 설정

### AWS 자격 증명 설정

DynamoGORM은 AWS SDK v2의 기본 설정을 사용합니다. 다음 중 하나의 방법으로 AWS 자격 증명을 설정하세요:

1. **AWS CLI 설정**
```bash
aws configure
```

2. **환경 변수**
```bash
export AWS_ACCESS_KEY_ID=your_access_key
export AWS_SECRET_ACCESS_KEY=your_secret_key
export AWS_REGION=ap-northeast-2
```

3. **IAM 역할** (EC2, Lambda 등에서 사용)

## 📖 사용 예시

### 1. 기본 구조체 정의

```go
package main

import (
    "fmt"
    "log"
    "github.com/zkfmapf123/dynamoGORM"
)

// 사용자 구조체 정의
type User struct {
    ID       string `json:"id"`        // Primary Key
    Name     string `json:"name"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
    IsActive bool   `json:"is_active"`
}

// 게시물 구조체 정의
type Post struct {
    ID        string   `json:"id"`         // Primary Key
    Title     string   `json:"title"`
    Content   string   `json:"content"`
    Tags      []string `json:"tags"`
    CreatedAt string   `json:"created_at"`
}
```

### 2. 데이터 삽입 (Insert)

```go
func main() {
    // 테이블 파라미터 설정
    tableParams := dynamodbgo.TableParmas{
        tableName:   "users",
        primarykey:  "id",
        billingMode: true, // On-Demand 모드
    }

    // 사용자 데이터 준비
    userData := map[string]any{
        "id":        "user123",
        "name":      "홍길동",
        "email":     "hong@example.com",
        "age":       30,
        "is_active": true,
    }

    // 데이터 삽입
    err := dynamodbgo.Insert(tableParams, userData)
    if err != nil {
        log.Fatal("Insert failed:", err)
    }
    fmt.Println("사용자 데이터 삽입 완료!")
}
```

### 3. 데이터 조회 (Select)

#### Primary Key로 단일 조회

```go
func getUserByID() {
    // ID로 사용자 조회
    user, err := dynamodbgo.FindByUsePK[User]("users", "id", "user123")
    if err != nil {
        log.Fatal("Find failed:", err)
    }
    
    fmt.Printf("조회된 사용자: %+v\n", user)
}
```

#### 전체 데이터 조회

```go
func getAllUsers() {
    // 모든 사용자 조회
    users, err := dynamodbgo.SelectAll[User]("users")
    if err != nil {
        log.Fatal("SelectAll failed:", err)
    }
    
    fmt.Printf("전체 사용자 수: %d\n", len(users))
    for _, user := range users {
        fmt.Printf("- %s: %s (%s)\n", user.ID, user.Name, user.Email)
    }
}
```

### 4. 데이터 수정 (Update)

```go
func updateUser() {
    // 수정할 데이터 준비
    updates := map[string]any{
        "name":      "김철수",
        "age":       25,
        "is_active": false,
    }

    // 사용자 정보 수정 (Primary Key 제외)
    err := dynamodbgo.UpdatePartial("users", "id", "user123", updates)
    if err != nil {
        log.Fatal("Update failed:", err)
    }
    
    fmt.Println("사용자 정보 수정 완료!")
}
```

### 5. 데이터 삭제 (Delete)

```go
func deleteUser() {
    // 사용자 삭제
    err := dynamodbgo.Delete("users", "id", "user123")
    if err != nil {
        log.Fatal("Delete failed:", err)
    }
    
    fmt.Println("사용자 삭제 완료!")
}
```

## ⚠️ 주의사항

1. **Primary Key**: 모든 테이블은 문자열 타입의 Primary Key가 필요합니다.
2. **테이블 생성**: 테이블이 없을 경우 자동으로 생성되며, 최대 10번까지 재시도합니다.
3. **빌링 모드**: On-Demand 모드(`true`) 또는 Provisioned 모드(`false`)를 선택할 수 있습니다.
4. **AWS 설정**: 사용하기 전에 AWS 자격 증명이 올바르게 설정되어야 합니다.

## 🤝 기여하기

버그 리포트나 기능 제안은 GitHub Issues를 통해 해주세요.

## 📄 라이선스

이 프로젝트는 MIT 라이선스 하에 배포됩니다.
