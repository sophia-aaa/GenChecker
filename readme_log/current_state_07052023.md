This code <br>
func (h *Header) GetUnsafePointer(i int) unsafe.Pointer { return h.UnsafePointers()[i] } <br>
will be changed via <b><i> main.go </i></b> <br>
<i>From ast package</i> <br>
...<br>
7117  .  .  54: *ast.FuncDecl {<br>
7118  .  .  .  Recv: *ast.FieldList {<br> 
7119  .  .  .  .  Opening: dataset/getset.go:149:6 <br>
...<br>
7143  .  .  .  } <br>
7144  .  .  .  Name: *ast.Ident { <br>
7145  .  .  .  .  NamePos: dataset/getset.go:149:18 <br>
7146  .  .  .  .  Name: "GetUnsafePointer" <br>
7147  .  .  .  } <br>
...<br>
7216  .  .  .  .  .  .  .  .  Index: *ast.Ident { <br>
7217  .  .  .  .  .  .  .  .  .  NamePos: dataset/getset.go:149:88 <br>
7218  .  .  .  .  .  .  .  .  .  Name: "i" <br>
7219  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 7158) <br>
7220  .  .  .  .  .  .  .  .  } <br>
...<br>

<br>
<i>To struct</i> <br>
...<br>
function name is  GetUnsafePointer <br>
{FieldList [h]} <br>
{StarExpr [Header]} <br>
{FieldList [i int]} <br>
{FieldList -> SelectorExpr [unsafe Pointer]} <br>
{BlockStmt -> ReturnStmt -> IndexExpr -> CallExpr -> SelectorExpr [h UnsafePointers i]} <br>
...<br>

Command on the terminal should be: <br>
go run main.go - "xxx.go" <br>
e.g. <i>go run main.go - "dataset/getset.go"</i>

