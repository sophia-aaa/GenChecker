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


patternArraySet
function 이름 그대로

if path contains *ast.BlockStmt
beforeBlockStm = false
-. method 인지 체크, 리시버 이름 저장
path contains *ast.FieldList -> *ast.Field ?
그 다음 path contains *ast.StarExpr ?
h랑 h 구조 이름 저장 receiverParam

0. 파라미터에 int 들어가있는지 체크
   path contains *ast.FieldList -> *ast.Field
   astValue 길이 체크랑 astValue[1] int ?
   astValue 저장

1. Parameter에 Interface가 있는지
   before beforeBlockStm true ? (to check whether it is parameter) and path contains *ast.Field
   그 다음에 path contains *ast.InterfaceType?
   해당 path ident 저장 = InterfaceVar


2. Switch statement begin and check whether the switch statement handles a type
   path contains *ast.SwitchStmt -> *ast.CallExpr -> *ast.SelectorExpr
   Ident contains recevierParam ?
   Ident contains Kind ?

3. Case Clause check
   path contains *ast.CaseClause -> *ast.SelectorExpr
   Ident contains reflect ?
   astValue 저장 = caseVar

4. Set Value 저장
   path contains *ast.TypeAssertExpr
   astValue contains paramIdent ? and astValue contains caseClausIdent ?
   그 전 path contains *ast.AssignStmt
   Ident 저장 typeVar

5. function check
   path contains *ast.ExprStmt -> *ast.CallExpr -> *ast.SelectorExpr
   astValue contains recevierParam ? && astValue contains funcName ?
   astValue contains param ?
   astValue contains typeVar ?

   
Replace

patternArrayGet
function 이름 그대로


if path contains *ast.BlockStmt
beforeBlockStm = false
-. method 인지 체크, 리시버 이름 저장
path contains *ast.FieldList -> *ast.Field ?
그 다음 path contains *ast.StarExpr ?
h랑 h 구조 이름 저장 receiverParam

0. 파라미터에 int 들어가있는지 체크
   path contains *ast.FieldList -> *ast.Field
   astValue 길이 체크랑 astValue[1] int ?
   astValue 저장

1. Parameter에 Interface가 있는지
   path contains *ast.FieldList -> *ast.Field -> *ast.InterfaceType


2. Switch statement begin and check whether the switch statement handles a type
   path contains *ast.SwitchStmt -> *ast.CallExpr -> *ast.SelectorExpr
   Ident contains recevierParam ?
   Ident contains Kind ?

3. Case Clause check
   path contains *ast.CaseClause -> *ast.SelectorExpr
   Ident contains reflect ?


4. function check
   path contains *ast.ReturnStmt -> *ast.CallExpr -> *ast.SelectorExpr
   astValue contains recevierParam ? && astValue contains funcName ? && astValue contains param ?


Replace

patternArrayMemset

function 이름 그대로


if path contains *ast.BlockStmt
beforeBlockStm = false
-. method 인지 체크, 리시버 이름 저장
path contains *ast.FieldList -> *ast.Field ?
그 다음 path contains *ast.StarExpr ?
h랑 h 구조 이름 저장 receiverParam



0. Parameter에 Interface가 있는지
   beforeBlockStm true && path contains *ast.InterfaceType && length of astValue == 0
   그 전 path contains *ast.FieldList -> *ast.Field ?
   astValue = interfaceVal

1. Save the return value

2. Switch statement begin and check whether the switch statement handles a type
   path contains **ast.SwitchStmt -> *ast.SelectorExpr
   Ident contains recevierParam ?

3. Case Clause check
   path contains *ast.CaseClause && length of astValue == 1   
   Ident contains type ?
   Ident = caseVar

4. If there is a interface variable in the astValue?
   path contains *ast.TypeAssertExpr
   Ident contains interfaceVal ? && Ident contains caseVar ?
   그 전 path contains *ast.AssignStmt ?
   astValue[1] = setVar


5. If there is a receiver Parameter in Astvalue
   path contains *ast.CallExpr -> *ast.SelectorExpr
   Ident contains a receiver Parameter ? && contains Type?
   그 전 path contains *ast.AssignStmt?
   Ident = arrayVar

6. if there is Range Statement?
   path contains *ast.RangeStmt
   astValue contains arrayVar?
   astValue[0] = rangeVar

7. if there is array expression ?
   path contains *ast.IndexExpr
   astValue contains arrayVar ? && astValue contains rangeVar ? && astValue contains setVar ?

8. Check Return value
   path contains *ast.ReturnStmt && astValue == nil ?
   break



Replace


patternMemsetIter
function 이름 그대로


if path contains *ast.BlockStmt
beforeBlockStm = false
-. method 인지 체크, 리시버 이름 저장
path contains *ast.FieldList -> *ast.Field ?
그 다음 path contains *ast.StarExpr ?
h랑 h 구조 이름 저장 receiverParam



0. Parameter에 Interface가 있는지
   beforeBlockStm true && path contains *ast.InterfaceType && length of astValue == 0
   그 전 path contains *ast.FieldList -> *ast.Field ?
   astValue = interfaceVal

1. Parameter에 Iterator 가 있는지
   beforeBlockStm true && path contains *ast.FieldList -> *ast.Field && length of astValue > 1
   astValue[0] = iteratorVar

2. Save the return value and type
   beforeBlockStm true && path contains *ast.FieldList -> *ast.Field && length of astValue > 1
   astValue[0] = returnValue

3. Save int value for for-loop
   path contains *ast.GenDecl -> *ast.ValueSpec && length of astValue > 1 && astValue[1] == int
   astValue[0] = intVar

4. Switch statement begin and check whether the switch statement handles a type
   path contains **ast.SwitchStmt -> *ast.SelectorExpr
   Ident contains recevierParam ?

5. Case Clause check
   path contains *ast.CaseClause && length of astValue == 1   
   Ident contains type ?
   Ident = caseVar

6. If there is a interface variable in the astValue?
   path contains *ast.TypeAssertExpr
   Ident contains interfaceVal ? && Ident contains caseVar ?
   그 전 path contains *ast.AssignStmt ?
   astValue[1] = setVar


7. If there is a receiver Parameter in Astvalue
   path contains *ast.CallExpr -> *ast.SelectorExpr
   Ident contains a receiver Parameter ? && contains Type?
   그 전 path contains *ast.AssignStmt?
   Ident = arrayVar

6. if there is a for-loop ?
   path contains *ast.ForStmt -> *ast.AssignStmt && astValue contains intVar && astValue contains returnValue
   그 다음 astValue == {iteratorVar, "Next"}
   그 그 다음 astValue == {"==", returnValue, "nil"}
   그 그 그 다음 astValue == {"=", intVar, returnValue}
   그 그 그 그 다음 value == {iteratorVar, "Next"}

7. if there is array expression ?
   path contains *ast.IndexExpr
   astValue contains arrayVar ? && astValue contains rangeVar ? && astValue contains setVar ?


8.  If there is a return value in Astvalue
    path contains *ast.AssignStmt && astValue contains returnValue &&
    그 다음 path contains *ast.CallExpr && length of astValue > 1 && astValue contains returnValue
    astValue[0] = errorFunc

Replace

patternEq
function 이름 그대로


if path contains *ast.BlockStmt
beforeBlockStm = false
-. method 인지 체크, 리시버 이름 저장
path contains *ast.FieldList -> *ast.Field ?
그 다음 path contains *ast.StarExpr ?
h랑 h 구조 이름 저장 receiverParam



0. Parameter에 Interface가 있는지
   beforeBlockStm true && path contains *ast.InterfaceType && length of astValue == 0
   그 전 path contains *ast.FieldList -> *ast.Field ?
   astValue = interfaceVal

1. If there is a interface variable in the astValue?
   path contains *ast.TypeAssertExpr
   Ident contains interfaceVal ?
   그 전 path contains *ast.AssignStmt ? length
   astValue[1] = interfaceVal1

2. Case Clause check
   path contains *ast.CaseClause -> *ast.SelectorExpr && length of astValue > 1 && astValue[0] == reflect   
   astValue[1] = caseVar

3. If there is a receiver Value in the astValue and for-Loop ?
   path contains *ast.CallExpr -> *ast.SelectorExpr
   Ident contains interfaceVal ? && Ident contains caseVar ?
   그 전 path contains *ast.RangeStmt ? && len(astValue) > 1
   astValue[0] = intVar
   astValue[1] = receiverRangeVar

4. if there is compare statement
   path contains *ast.BinaryExpr
   astValue contains "==" || astValue contains "!="
   그 다음 path contains *ast.CallExpr -> *ast.SelectorExpr
   for Loop range astValue
   if astValue contains interfaceVal1 then next value not receiverRangeVar
   -> no parameterized methods

patternReduce

