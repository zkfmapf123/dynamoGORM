# DynamoDB Single Table 설계 가이드

## 목차

- [핵심 개념](#핵심-개념)
- [PK vs SK 구분법](#pk-vs-sk-구분법)
- [설계 프로세스](#설계-프로세스)
- [자주 헷갈리는 케이스](#자주-헷갈리는-케이스)
- [실전 예제](#실전-예제)
- [안티패턴](#안티패턴)

---

## 핵심 개념

### DynamoDB의 본질

DynamoDB = 초고속 Key-Value 스토어
- 복잡한 SQL 없음
- JOIN 없음
- 설계가 곧 쿼리

### Single Table Design

```
MySQL:     여러 테이블 → JOIN으로 연결
DynamoDB:  하나의 테이블 → PK/SK로 구분
```

### 가장 중요한 원칙

```
PK = 큰 그룹 (폴더)
SK = 그 안의 세부 항목들 (파일)
```

---

## PK vs SK 구분법

### 1단계: 접근 패턴 분석

질문: "무엇을 기준으로 조회하나?"

```
예시: "특정 작성자의 모든 글"
→ "작성자"를 기준으로 조회
→ PK에 작성자 정보

예시: "특정 룸의 모든 메시지"
→ "룸"을 기준으로 조회
→ PK에 룸 정보
```

### 2단계: PK/SK 결정

패턴: "A의 모든 B를 조회"
→ PK: A
→ SK: B

#### 실전 예시

```
"유저의 모든 주문"
→ PK: USER#123
→ SK: ORDER#001, ORDER#002, ...

"룸의 모든 메시지"
→ PK: ROOM#abc
→ SK: MSG#001, MSG#002, ...

"개발팀의 모든 멤버"
→ PK: TEAM#DEV
→ SK: USER#1, USER#2, ...
```

---

## 설계 프로세스

### Step 1: 접근 패턴 나열

```
1. 특정 유저 정보 조회
2. 유저의 모든 주문 조회
3. 특정 주문의 상세 정보 조회
4. 특정 상품이 포함된 모든 주문 조회
```

### Step 2: 엔티티 파악

- User (유저)
- Order (주문)
- Product (상품)

### Step 3: PK/SK 설계

```
패턴 1, 2 → 유저 기준
PK: USER#123
├─ SK: #PROFILE      (유저 정보)
└─ SK: ORDER#001     (주문1)

패턴 3 → 주문 기준
PK: ORDER#001
├─ SK: #METADATA     (주문 정보)
└─ SK: ITEM#PROD-456 (주문 상품)

패턴 4 → 상품 기준 (GSI 또는 역방향 저장)
PK: PRODUCT#456
└─ SK: ORDER#001     (이 상품이 포함된 주문)
```

---

## 자주 헷갈리는 케이스

### 케이스 1: PK에 뭘 넣어야 할까?

#### ❌ 잘못된 생각

```go
PK: "USER"           // 모든 유저?
PK: "USER#DEV"       // 개발팀 유저들?
PK: "EMAIL"          // 이메일?
```

#### ✅ 올바른 설계

```go
PK: "USER#123"       // 특정 유저 (개별 ID)
PK: "TEAM#DEV"       // 특정 팀 (그룹 ID)
PK: "ORDER#456"      // 특정 주문 (개별 ID)
```

**핵심: PK = {엔티티타입}#{고유ID}**

---

### 케이스 2: # 기호는 어디에?

#### ❌ 잘못된 이해

```go
PK: "USER"           SK: "DEV#123"
PK: "USER#DEV"       SK: "PROFILE"
```

#### ✅ 올바른 설계

```go
// # = 구분자
PK: "USER#123"       SK: "#PROFILE"
     ↑    ↑              ↑
   타입  ID           메타데이터

PK: "TEAM#DEV"       SK: "USER#123"
     ↑    ↑              ↑     ↑
   타입  ID            타입   ID
```

**규칙:**
- PK: `{타입}#{ID}`
- SK: `{타입}#{추가정보}` 또는 `#{메타데이터}`

---

### 케이스 3: 특정 유저는 어떻게 찾아?

#### 질문

```
PK: "USER#123"으로 저장했는데,
나중에 이 유저를 어떻게 찾지?
Scan 밖에 못하지 않아?
```

#### 해답: 3가지 방법

**방법 1: ID를 알고 있는 경우 (가장 일반적)**

```go
// 로그인 시 userID 획득
// 세션/JWT에 저장: userID = "123"

// 이후 조회
GetItem(PK: "USER#123", SK: "#PROFILE")
```

**방법 2: Email로 찾기 (GSI)**

```go
// 테이블 설계
{
    PK: "USER#123",
    SK: "#PROFILE",
    Email: "kim@email.com",
    
    // GSI
    GSI1PK: "kim@email.com",
    GSI1SK: "USER"
}

// 조회
Query(IndexName: "GSI1", GSI1PK: "kim@email.com")
→ PK: "USER#123" 획득
```

**방법 3: Email을 PK로 (간단하지만 제약)**

```go
PK: "kim@email.com"
SK: "#PROFILE"

// 장점: GSI 불필요, 간단
// 단점: Email 변경 불가, 확장성 제한
```

---

### 케이스 4: 양방향 조회가 필요하면?

#### 시나리오

```
1. 특정 유저 정보 조회
2. 특정 팀의 모든 멤버 조회
```

#### 해결: 데이터 중복 저장

```go
// 유저 정보 (패턴 1용)
{
    PK: "USER#1",
    SK: "#PROFILE",
    Name: "tom",
    Team: "DEV"
}

// 팀별 멤버 (패턴 2용)
{
    PK: "TEAM#DEV",
    SK: "USER#1",
    Name: "tom"
}
```

**핵심: DynamoDB에서 중복 저장 = 정상**

---

### 케이스 5: SK는 WHERE절이다

#### SQL 사고방식

```sql
SELECT * FROM users 
WHERE department = 'DEV' 
  AND age BETWEEN 20 AND 30
ORDER BY age;
```

#### DynamoDB 설계

```go
PK: "DEPT#DEV"                    // WHERE department = 'DEV'
SK: "AGE#0025#USER#123"           // WHERE age + ORDER BY age
     ↑     ↑
   정렬기준 정렬값(제로패딩)
```

**핵심: 조회 조건 = SK 설계**

---

### 케이스 6: 정렬이 필요하면?

#### 요구사항

```
모든 유저를 나이순으로 조회
```

#### 설계

```go
// PK 통일 + SK에 나이
PK: "USERS"             SK: "AGE#0025#USER#1"
PK: "USERS"             SK: "AGE#0030#USER#2"
PK: "USERS"             SK: "AGE#0035#USER#3"

// 쿼리
Query(PK: "USERS")
// 결과: 자동으로 나이순 정렬 (SK 기준)
```

**중요: 숫자는 제로패딩 필수!**

```
❌ "AGE#5", "AGE#25", "AGE#100"  → "100" < "25" < "5" (잘못됨)
✅ "AGE#005", "AGE#025", "AGE#100" → 올바른 정렬
```

---

## 실전 예제

### 예제 1: 블로그 시스템

#### 접근 패턴

```
1. 특정 작성자의 모든 글
2. 특정 글의 모든 댓글 (최신순)
3. 특정 카테고리의 모든 글
```

#### 설계

```go
// 작성자의 글 (패턴 1)
PK: "AUTHOR#1"          SK: "POST#001"
PK: "AUTHOR#1"          SK: "POST#002"

// 글 정보 + 댓글 (패턴 2)
PK: "POST#001"          SK: "#METADATA"
PK: "POST#001"          SK: "COMMENT#2024-01-10T10:00#001"
PK: "POST#001"          SK: "COMMENT#2024-01-10T11:00#002"

// 카테고리별 글 (패턴 3)
PK: "CATEGORY#tech"     SK: "POST#001"
PK: "CATEGORY#tech"     SK: "POST#003"
```

#### 쿼리

```go
// 패턴 1
Query(PK: "AUTHOR#1", SK begins_with "POST#")

// 패턴 2
Query(
    PK: "POST#001",
    SK begins_with "COMMENT#",
    ScanIndexForward: false  // 최신순
)

// 패턴 3
Query(PK: "CATEGORY#tech", SK begins_with "POST#")
```

---

### 예제 2: 팀 관리 시스템

#### 접근 패턴

```
1. 특정 유저 정보 조회
2. 개발팀 멤버 조회
3. 프로덕트팀 멤버 조회
```

#### 설계

```go
// 유저 정보 (패턴 1)
{
    PK: "USER#1",
    SK: "#PROFILE",
    Name: "tom",
    Age: 32,
    Team: "DEV"
}

// 팀별 멤버 (패턴 2, 3)
{
    PK: "TEAM#DEV",
    SK: "USER#1",
    Name: "tom",
    Age: 32
}

{
    PK: "TEAM#DEV",
    SK: "USER#2",
    Name: "jane",
    Age: 28
}

{
    PK: "TEAM#PRODUCT",
    SK: "USER#3",
    Name: "mike",
    Age: 30
}
```

#### 코드

```go
// 저장: 각 유저마다 2개 아이템
// 1. 유저 정보
client.Insert({
    PK: "USER#1",
    SK: "#PROFILE",
    Name: "tom",
    Age: 32
})

// 2. 팀별 멤버
client.Insert({
    PK: "TEAM#DEV",
    SK: "USER#1",
    Name: "tom",
    Age: 32
})
```

#### 쿼리

```go
// 패턴 1
GetItem(PK: "USER#1", SK: "#PROFILE")

// 패턴 2
Query(PK: "TEAM#DEV")

// 패턴 3
Query(PK: "TEAM#PRODUCT")
```

---

## 안티패턴

### ❌ 안티패턴 1: PK가 모호

```go
PK: "USER"              // 누구? 모든 유저?
PK: "DATA"              // 무슨 데이터?
PK: "INFO"              // 뭔 정보?
```

### ✅ 올바른 패턴

```go
PK: "USER#123"          // 명확한 유저
PK: "ORDER#456"         // 명확한 주문
PK: "PRODUCT#789"       // 명확한 상품
```

---

### ❌ 안티패턴 2: SK가 비어있음

```go
{
    PK: "USER#1",
    SK: ""               // ❌ 빈 값
}
```

### ✅ 올바른 패턴

```go
{
    PK: "USER#1",
    SK: "#PROFILE"       // ✅ 명확한 값
}
```

---

### ❌ 안티패턴 3: 일반 컬럼으로 조회 시도

```go
// Age로 조회?
Query(Age > 20)          // ❌ 불가능!
```

### ✅ 올바른 패턴

```go
// SK에 Age 포함
PK: "USERS"
SK: "AGE#0025#USER#1"

Query(
    PK: "USERS",
    SK BETWEEN "AGE#0020" AND "AGE#0030"
)
```

---

### ❌ 안티패턴 4: PK가 자주 바뀜

```go
PK: "kim@email.com"      // Email 변경 시 문제!
```

### ✅ 올바른 패턴

```go
PK: "USER#123"           // 변경 불가능한 ID
Email: "kim@email.com"   // 일반 속성 (변경 가능)
GSI1PK: "kim@email.com"  // GSI로 조회
```

---

## 핵심 체크리스트

### 설계 전

- [ ] 모든 접근 패턴 나열했는가?
- [ ] 어떤 기준으로 조회하는가?
- [ ] 정렬이 필요한가?

### 설계 시

- [ ] PK = {타입}#{고유ID} 형식인가?
- [ ] SK에 조회 조건이 포함되었는가?
- [ ] 양방향 조회를 위한 중복 저장이 필요한가?
- [ ] 숫자 정렬 시 제로패딩 했는가?

### 설계 후

- [ ] 모든 접근 패턴이 Query/GetItem으로 가능한가?
- [ ] Scan이 필요한 부분은 없는가?
- [ ] GSI가 꼭 필요한가? (비용 고려)

---

## 마지막 팁

### DynamoDB는 MySQL과 다르다

```
MySQL:    테이블 만들고 → 나중에 쿼리 작성
DynamoDB: 쿼리 먼저 정의 → 테이블 설계
```

### 설계 = 쿼리

```
PK/SK 설계 = WHERE절 + ORDER BY절
잘못된 설계 = 느린 쿼리 (Scan)
```

### 중복 저장은 정상

```
MySQL:    중복 = 나쁜 설계
DynamoDB: 중복 = 좋은 설계 (성능 최적화)
```

### 접근 패턴이 모든 것

```
접근 패턴 명확 → 설계 쉬움
접근 패턴 불명확 → 설계 불가능
```

---

## 설계 원칙을 기억하세요

```
PK = 큰 그룹
SK = 세부 항목 + WHERE절
조회 패턴 = 설계의 시작
중복 저장 = 정상
```
