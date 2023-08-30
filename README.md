# GenChecker

This tool <i>GenChecker</i> was developed for the thesis "Code Transformation for Safety: Leveraging Generics to Refactor Unsafe Package Usages in Go" to replace the usage of unsafe package by generics in Go.
<br>This project was developed using the ast, astutil packages in go.
Type
```
go run . - "<path>/<filename.go>"
```
to run the tool.

The tool is divided into three steps.

First, the tool takes a single file as input and changes it into the basic structure of the tool using Go's ast package 
(see <i>buildAST.go</i>). Every running the tool will be produced a text file that the input file converts into the base structure of this tool.<br>
The text files are for debugging or analyse how the input file converted into the AST tree with Go's AST package.

Then, after the file is transformed into a base structure, the tool checks for reused code across multiple data types to match Go's generics recommendation usage. <br>
See <i>checkReuseFunc.go</i> for functions and <i>checkReuseFunc.go</i> for switch statements.<br>
Reference: https://go.dev/blog/when-generics
<br><br>
Finally, if any of the reused code matches the pattern, the file is refactored. 
<br>It will be suffixed with _replaced.

For patterns, we extracted patterns from the files labeled generics in the paper <i>"Uncovering the Hidden Dangers: Finding Unsafe Go Code in the Wild"</i> 
and <i>"UNGOML: Automated Classification of unsafe Usages in Go".</i> 
We extracted 10 patterns from 4 files, and 3 of them say that they cannot be replaced with generics. See pattern1.go, pattern2.go, pattern3.go, and pattern4.go for detailed patterns. 
<br>At the top of the file, I've outlined what code will be replaced.

You can use the file <i>spewConverter.go</i> to output the ast.Node of the files. <br>
This file relies on the package spew.

Due to the limitations of the ast and astutil packages, the tool can only refactor if there is code written in the same way as provided by the first dataset. 
<br>However, you can use the tool if the code contains reused implementations for multiple data types.

I strongly recommend you to adapt your keyword for your project to add your keyword according to your project in the const.go. Then the tool will be recognised as a list in the <i>const.go</i>

There may be gaps in the tool. Please fork and post an issue with any improvements or suggestions.
Thank you for your support and have a great day!
