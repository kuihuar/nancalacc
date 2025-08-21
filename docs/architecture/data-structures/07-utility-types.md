# 工具类 (Utility) - Utility Types

## 概述
工具类提供通用的辅助功能，包括字符串处理、时间处理、加密解密、数据转换等常用操作。

## 分类

### 7.1 字符串工具
```go
// 字符串工具类
type StringUtils struct{}

// 生成UUID
func (su *StringUtils) GenerateUUID() string {
    return uuid.New().String()
}

// 生成随机字符串
func (su *StringUtils) GenerateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[rand.Intn(len(charset))]
    }
    return string(b)
}

// 驼峰转下划线
func (su *StringUtils) CamelToSnake(s string) string {
    var result strings.Builder
    for i, r := range s {
        if i > 0 && unicode.IsUpper(r) {
            result.WriteRune('_')
        }
        result.WriteRune(unicode.ToLower(r))
    }
    return result.String()
}

// 下划线转驼峰
func (su *StringUtils) SnakeToCamel(s string) string {
    var result strings.Builder
    capitalize := true
    for _, r := range s {
        if r == '_' {
            capitalize = true
        } else {
            if capitalize {
                result.WriteRune(unicode.ToUpper(r))
                capitalize = false
            } else {
                result.WriteRune(r)
            }
        }
    }
    return result.String()
}

// 截断字符串
func (su *StringUtils) Truncate(s string, maxLength int) string {
    if len(s) <= maxLength {
        return s
    }
    return s[:maxLength-3] + "..."
}

// 检查字符串是否为空
func (su *StringUtils) IsEmpty(s string) bool {
    return strings.TrimSpace(s) == ""
}

// 检查字符串是否为有效邮箱
func (su *StringUtils) IsValidEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return emailRegex.MatchString(email)
}

// 检查字符串是否为有效手机号
func (su *StringUtils) IsValidPhone(phone string) bool {
    phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
    return phoneRegex.MatchString(phone)
}
```

### 7.2 时间工具
```go
// 时间工具类
type TimeUtils struct{}

// 获取当前时间戳
func (tu *TimeUtils) NowTimestamp() int64 {
    return time.Now().Unix()
}

// 获取当前毫秒时间戳
func (tu *TimeUtils) NowTimestampMillis() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}

// 格式化时间
func (tu *TimeUtils) FormatTime(t time.Time, layout string) string {
    return t.Format(layout)
}

// 解析时间字符串
func (tu *TimeUtils) ParseTime(timeStr, layout string) (time.Time, error) {
    return time.Parse(layout, timeStr)
}

// 获取时间差
func (tu *TimeUtils) TimeDiff(start, end time.Time) time.Duration {
    return end.Sub(start)
}

// 检查是否为同一天
func (tu *TimeUtils) IsSameDay(t1, t2 time.Time) bool {
    y1, m1, d1 := t1.Date()
    y2, m2, d2 := t2.Date()
    return y1 == y2 && m1 == m2 && d1 == d2
}

// 获取本周开始时间
func (tu *TimeUtils) WeekStart(t time.Time) time.Time {
    weekday := t.Weekday()
    if weekday == time.Sunday {
        weekday = 7
    } else {
        weekday--
    }
    return t.AddDate(0, 0, -int(weekday))
}

// 获取本月开始时间
func (tu *TimeUtils) MonthStart(t time.Time) time.Time {
    return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// 获取本年开始时间
func (tu *TimeUtils) YearStart(t time.Time) time.Time {
    return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}
```

### 7.3 加密工具
```go
// 加密工具类
type CryptoUtils struct {
    secretKey []byte
}

// 创建加密工具实例
func NewCryptoUtils(secretKey string) *CryptoUtils {
    return &CryptoUtils{
        secretKey: []byte(secretKey),
    }
}

// AES加密
func (cu *CryptoUtils) AESEncrypt(plaintext string) (string, error) {
    block, err := aes.NewCipher(cu.secretKey)
    if err != nil {
        return "", err
    }
    
    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }
    
    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))
    
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AES解密
func (cu *CryptoUtils) AESDecrypt(ciphertext string) (string, error) {
    data, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }
    
    block, err := aes.NewCipher(cu.secretKey)
    if err != nil {
        return "", err
    }
    
    if len(data) < aes.BlockSize {
        return "", fmt.Errorf("ciphertext too short")
    }
    
    iv := data[:aes.BlockSize]
    data = data[aes.BlockSize:]
    
    stream := cipher.NewCFBDecrypter(block, iv)
    stream.XORKeyStream(data, data)
    
    return string(data), nil
}

// MD5哈希
func (cu *CryptoUtils) MD5Hash(data string) string {
    hash := md5.Sum([]byte(data))
    return hex.EncodeToString(hash[:])
}

// SHA256哈希
func (cu *CryptoUtils) SHA256Hash(data string) string {
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}

// 生成随机盐
func (cu *CryptoUtils) GenerateSalt(length int) string {
    salt := make([]byte, length)
    rand.Read(salt)
    return base64.StdEncoding.EncodeToString(salt)
}

// 密码哈希
func (cu *CryptoUtils) HashPassword(password, salt string) string {
    return cu.SHA256Hash(password + salt)
}
```

### 7.4 数据转换工具
```go
// 数据转换工具类
type ConvertUtils struct{}

// 字符串转整数
func (cu *ConvertUtils) StringToInt(s string, defaultValue int) int {
    if i, err := strconv.Atoi(s); err == nil {
        return i
    }
    return defaultValue
}

// 字符串转浮点数
func (cu *ConvertUtils) StringToFloat(s string, defaultValue float64) float64 {
    if f, err := strconv.ParseFloat(s, 64); err == nil {
        return f
    }
    return defaultValue
}

// 字符串转布尔值
func (cu *ConvertUtils) StringToBool(s string, defaultValue bool) bool {
    if b, err := strconv.ParseBool(s); err == nil {
        return b
    }
    return defaultValue
}

// 接口转字符串
func (cu *ConvertUtils) InterfaceToString(v interface{}) string {
    if v == nil {
        return ""
    }
    switch val := v.(type) {
    case string:
        return val
    case int, int8, int16, int32, int64:
        return strconv.FormatInt(reflect.ValueOf(val).Int(), 10)
    case uint, uint8, uint16, uint32, uint64:
        return strconv.FormatUint(reflect.ValueOf(val).Uint(), 10)
    case float32, float64:
        return strconv.FormatFloat(reflect.ValueOf(val).Float(), 'f', -1, 64)
    case bool:
        return strconv.FormatBool(val)
    default:
        return fmt.Sprintf("%v", val)
    }
}

// 结构体转Map
func (cu *ConvertUtils) StructToMap(obj interface{}) map[string]interface{} {
    result := make(map[string]interface{})
    v := reflect.ValueOf(obj)
    t := v.Type()
    
    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fieldType := t.Field(i)
        
        // 获取JSON标签
        jsonTag := fieldType.Tag.Get("json")
        if jsonTag == "" || jsonTag == "-" {
            continue
        }
        
        // 处理omitempty
        if strings.Contains(jsonTag, ",") {
            jsonTag = strings.Split(jsonTag, ",")[0]
        }
        
        result[jsonTag] = field.Interface()
    }
    
    return result
}

// Map转结构体
func (cu *ConvertUtils) MapToStruct(data map[string]interface{}, obj interface{}) error {
    v := reflect.ValueOf(obj)
    if v.Kind() != reflect.Ptr {
        return fmt.Errorf("obj must be a pointer")
    }
    
    v = v.Elem()
    t := v.Type()
    
    for i := 0; i < v.NumField(); i++ {
        field := v.Field(i)
        fieldType := t.Field(i)
        
        jsonTag := fieldType.Tag.Get("json")
        if jsonTag == "" || jsonTag == "-" {
            continue
        }
        
        if strings.Contains(jsonTag, ",") {
            jsonTag = strings.Split(jsonTag, ",")[0]
        }
        
        if value, exists := data[jsonTag]; exists {
            field.Set(reflect.ValueOf(value))
        }
    }
    
    return nil
}
```

### 7.5 集合工具
```go
// 集合工具类
type CollectionUtils struct{}

// 切片去重
func (cu *CollectionUtils) RemoveDuplicates(slice []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    
    for _, item := range slice {
        if !seen[item] {
            seen[item] = true
            result = append(result, item)
        }
    }
    
    return result
}

// 切片交集
func (cu *CollectionUtils) Intersection(slice1, slice2 []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    
    for _, item := range slice1 {
        seen[item] = true
    }
    
    for _, item := range slice2 {
        if seen[item] {
            result = append(result, item)
        }
    }
    
    return result
}

// 切片并集
func (cu *CollectionUtils) Union(slice1, slice2 []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    
    for _, item := range slice1 {
        if !seen[item] {
            seen[item] = true
            result = append(result, item)
        }
    }
    
    for _, item := range slice2 {
        if !seen[item] {
            seen[item] = true
            result = append(result, item)
        }
    }
    
    return result
}

// 切片差集
func (cu *CollectionUtils) Difference(slice1, slice2 []string) []string {
    seen := make(map[string]bool)
    result := []string{}
    
    for _, item := range slice2 {
        seen[item] = true
    }
    
    for _, item := range slice1 {
        if !seen[item] {
            result = append(result, item)
        }
    }
    
    return result
}

// 切片分块
func (cu *CollectionUtils) Chunk(slice []string, size int) [][]string {
    var chunks [][]string
    for i := 0; i < len(slice); i += size {
        end := i + size
        if end > len(slice) {
            end = len(slice)
        }
        chunks = append(chunks, slice[i:end])
    }
    return chunks
}
```

## 特点

1. **通用性**: 提供通用的工具函数
2. **可复用**: 可在多个模块中重复使用
3. **无状态**: 工具类通常是无状态的
4. **纯函数**: 相同的输入总是产生相同的输出
5. **高性能**: 经过优化的实现

## 使用场景

- 数据处理
- 格式转换
- 验证检查
- 加密解密
- 字符串处理

## 最佳实践

1. **单一职责**: 每个工具类专注于特定功能
2. **静态方法**: 使用静态方法或函数
3. **错误处理**: 提供适当的错误处理
4. **性能优化**: 考虑性能影响
5. **测试覆盖**: 提供完整的单元测试 