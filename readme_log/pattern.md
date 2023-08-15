patternObject process

function 이름 앞에 대문자인지 아닌지 체크
대문자라면 ObjectSlice
소문자라면 objectSlice

1. method 인지 체크, 리시버 이름 저장
   path contains *ast.FieldList -> *ast.Field ?
   h 이름 저장

2. return type slice type인지 체크
   path contains *ast.Field -> *ast.ArrayType
   ident contains type ?  -> 이 처리는 나중에
   true 저장
   false 저장 but interface flag true

3. return value return (*(*[]bool) ... 이렇게 생겼는지 체크
   path contains *ast.ReturnStmt -> *ast.SliceExpr -> *ast.ParenExpr -> *ast.StarExpr -> *ast.CallExpr -> *ast.ParenExpr -> *ast.StarExpr -> *ast.ArrayType
   ident same as return type ?


4. unsafe Pointer check
   path contains *ast.SelectorExpr
   ident contains unsafe Pointer ?

5. UnaryExpr check
   path contains *ast.UnaryExpr
   ident == & ?

6. 리시버에 있는 필드를 부르는게 맞는지 확인
   path contains *ast.SelectorExpr
   ident 리시버 이름이랑 끝에서 두번째 두번쨰 이름이 같은지
   맨 끝에 있는 이름이 리스트 이름


Replace


patternSet process

function 이름 앞에 대문자인지 아닌지 체크 && 이름에 set 들어가있는지 체크
대문자라면 Set
소문자라면 set


1. method 인지 체크, 리시버 이름 저장
   path contains *ast.FieldList -> *ast.Field ?
   그 다음 path contains *ast.StarExpr ? 이거 중요!
   h랑 h 구조 이름 저장

2. 파라미터에 int 들어가있는지 체크
   path contains *ast.FieldList -> *ast.Field
   astValue 길이 체크랑 astValue 저장
   astValue[1] int ?


3. return value 저장
   path contains *ast.Field
   ident same as return type ?


4. Assignment check
   path contains  *ast.AssignStmt
   ident == =
   4.1 IndexExpr check
   path contains *ast.IndexExpr -> *ast.CallExpr -> *ast.SelectorExpr
   astValue에 parameter, receiver parameter 그리고 return value 저장되어있는지 체크


Replace

function 이름 앞에 대문자인지 아닌지 체크 && 이름에 set 들어가있는지 체크
대문자라면 Set
소문자라면 set


1. method 인지 체크, 리시버 이름 저장
   path contains *ast.FieldList -> *ast.Field ?
   그 다음 path contains *ast.StarExpr ? 이거 중요!
   h랑 h 구조 이름 저장

2. 파라미터에 int 들어가있는지 체크
   path contains *ast.FieldList -> *ast.Field
   astValue 길이 체크랑 astValue 저장
   astValue[1] int ?


3. return value 저장
   path contains *ast.FieldList -> *ast.Field
   ident 저장


4. return statement check
   path contains  *ast.ReturnStmt -> *ast.IndexExpr -> *ast.CallExpr -> *ast.SelectorExpr
   astValue에 parameter, receiver parameter 저장되어있는지 체크


Replace

patternGet process