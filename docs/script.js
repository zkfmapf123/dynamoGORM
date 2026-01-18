// ========================================
// Code Examples Data
// ========================================
const codeExamples = {
    client: `package main

import (
    "context"

    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
    gdrm "github.com/zkfmapf123/gdrm"
)

func main() {
    // AWS 설정 로드
    cfg, _ := config.LoadDefaultConfig(context.Background())
    dynamoClient := dynamodb.NewFromConfig(cfg)

    // GDRM 클라이언트 생성
    client := gdrm.NewDDB(dynamoClient)

    // 테이블 설정 추가
    client.AddTable("my_table", gdrm.DDBTableParams{
        IsCreate:        true,
        IsPK:            true,
        PkAttributeType: types.ScalarAttributeTypeS,
        IsSK:            true,
        SkAttributeType: types.ScalarAttributeTypeS,
    })

    // 테이블 생성 시작
    client.Start(true)
}`,

    insert: `type User struct {
    PK   string \`dynamodbav:"PK"\`
    SK   string \`dynamodbav:"SK"\`
    Name string \`dynamodbav:"Name"\`
    Age  int    \`dynamodbav:"Age"\`
}

// 단건 삽입
err := client.Insert("my_table", User{
    PK:   "USER#123",
    SK:   "#PROFILE",
    Name: "tom",
    Age:  32,
})

// 배치 삽입 (25개씩 자동 분할)
users := []any{
    User{PK: "USER#1", SK: "#PROFILE", Name: "tom", Age: 32},
    User{PK: "USER#2", SK: "#PROFILE", Name: "jane", Age: 28},
    User{PK: "USER#3", SK: "#PROFILE", Name: "mike", Age: 30},
}

err = client.InsertBatch("my_table", users)`,

    select: `// 단건 조회
item, err := client.FindByKey("my_table", "USER#123", "#PROFILE")
if err != nil {
    log.Fatal(err)
}

// Expression을 사용한 조회
items, err := client.FindByKeyUseExpression(
    "my_table",
    "TEAM#DEV",
    "",
    100,  // limit
    gdrm.RangeParams{
        KeyConditionExpression: "PK = :pk",
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":pk": &types.AttributeValueMemberS{Value: "TEAM#DEV"},
        },
    },
)

// begins_with 사용
items, err = client.FindByKeyUseExpression(
    "my_table",
    "USER#123",
    "ORDER#",
    50,
    gdrm.RangeParams{
        KeyConditionExpression: "PK = :pk AND begins_with(SK, :sk)",
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":pk": &types.AttributeValueMemberS{Value: "USER#123"},
            ":sk": &types.AttributeValueMemberS{Value: "ORDER#"},
        },
    },
)`,

    marshal: `// 단건 결과 변환
item, _ := client.FindByKey("my_table", "USER#123", "#PROFILE")

// 제네릭을 사용한 타입 변환
user := gdrm.MarshalMap[User](item)

fmt.Println(user.Name)  // "tom"
fmt.Println(user.Age)   // 32

// 복수 결과 변환
items, _ := client.FindByKeyUseExpression(...)

users := gdrm.MarshalMaps[User](items)

for _, u := range users {
    fmt.Printf("%s: %d세\\n", u.Name, u.Age)
}`
};

// ========================================
// DOM Elements
// ========================================
const codeDisplay = document.getElementById('code-display');
const tabButtons = document.querySelectorAll('.tab-btn');

// ========================================
// Tab Switching
// ========================================
function switchTab(tabName) {
    // Update active tab
    tabButtons.forEach(btn => {
        btn.classList.toggle('active', btn.dataset.tab === tabName);
    });

    // Update code display with animation
    codeDisplay.style.opacity = '0';
    
    setTimeout(() => {
        codeDisplay.textContent = codeExamples[tabName];
        highlightCode();
        codeDisplay.style.opacity = '1';
    }, 150);
}

// Initialize tab click handlers
tabButtons.forEach(btn => {
    btn.addEventListener('click', () => {
        switchTab(btn.dataset.tab);
    });
});

// ========================================
// Code Copy Function
// ========================================
function copyCode() {
    const code = codeDisplay.textContent;
    navigator.clipboard.writeText(code).then(() => {
        const copyBtn = document.querySelector('.copy-btn');
        const originalText = copyBtn.textContent;
        
        copyBtn.textContent = '복사됨!';
        copyBtn.classList.add('copied');
        
        setTimeout(() => {
            copyBtn.textContent = originalText;
            copyBtn.classList.remove('copied');
        }, 2000);
    });
}

// ========================================
// Basic Syntax Highlighting
// ========================================
function highlightCode() {
    let code = codeDisplay.innerHTML;
    
    // Keywords
    const keywords = ['package', 'import', 'func', 'type', 'struct', 'return', 'if', 'for', 'range', 'var', 'const', 'err', 'nil', 'true', 'false'];
    keywords.forEach(keyword => {
        const regex = new RegExp(`\\b(${keyword})\\b`, 'g');
        code = code.replace(regex, '<span class="keyword">$1</span>');
    });
    
    // Strings
    code = code.replace(/(["'`])(?:(?!\1)[^\\]|\\.)*\1/g, '<span class="string">$&</span>');
    
    // Comments
    code = code.replace(/(\/\/.*)/g, '<span class="comment">$1</span>');
    
    // Types
    const types = ['string', 'int', 'bool', 'error', 'any', 'context', 'User'];
    types.forEach(type => {
        const regex = new RegExp(`\\b(${type})\\b`, 'g');
        code = code.replace(regex, '<span class="type">$1</span>');
    });
    
    codeDisplay.innerHTML = code;
}

// ========================================
// Smooth Scroll for Navigation
// ========================================
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function(e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            const offset = 80; // navbar height
            const targetPosition = target.getBoundingClientRect().top + window.pageYOffset - offset;
            
            window.scrollTo({
                top: targetPosition,
                behavior: 'smooth'
            });
        }
    });
});

// ========================================
// Navbar Background on Scroll
// ========================================
let lastScrollY = window.scrollY;

window.addEventListener('scroll', () => {
    const navbar = document.querySelector('.navbar');
    
    if (window.scrollY > 100) {
        navbar.style.background = 'rgba(10, 10, 11, 0.95)';
    } else {
        navbar.style.background = 'rgba(10, 10, 11, 0.8)';
    }
    
    lastScrollY = window.scrollY;
});

// ========================================
// Intersection Observer for Animations
// ========================================
const observerOptions = {
    root: null,
    rootMargin: '0px',
    threshold: 0.1
};

const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
        if (entry.isIntersecting) {
            entry.target.style.opacity = '1';
            entry.target.style.transform = 'translateY(0)';
        }
    });
}, observerOptions);

// Observe elements for animation
document.querySelectorAll('.func-doc, .design-card, .install-step').forEach(el => {
    el.style.opacity = '0';
    el.style.transform = 'translateY(20px)';
    el.style.transition = 'opacity 0.6s ease, transform 0.6s ease';
    observer.observe(el);
});

// ========================================
// Active TOC Link Highlighting
// ========================================
const tocLinks = document.querySelectorAll('.toc-group a');
const sections = document.querySelectorAll('.func-doc');

window.addEventListener('scroll', () => {
    let current = '';
    
    sections.forEach(section => {
        const sectionTop = section.offsetTop;
        const sectionHeight = section.clientHeight;
        
        if (window.scrollY >= sectionTop - 150) {
            current = section.getAttribute('id');
        }
    });
    
    tocLinks.forEach(link => {
        link.style.color = '';
        link.style.background = '';
        
        if (link.getAttribute('href') === `#${current}`) {
            link.style.color = 'var(--color-primary)';
            link.style.background = 'rgba(245, 158, 11, 0.1)';
        }
    });
});

// ========================================
// Initialize
// ========================================
document.addEventListener('DOMContentLoaded', () => {
    // Set initial code example
    switchTab('client');
    
    // Add transition to code display
    codeDisplay.style.transition = 'opacity 0.15s ease';
});

