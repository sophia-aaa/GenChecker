ABSTRACT - Today, Go is a widely used programming language despite its relatively short history. 
Due to its strong type system and simple language structure, Go provides reliable and compatible production that functions with different operating systems. 
The language itself has memory management using a garbage collector, so it is considered a memory safe language (US National Security Agency.
2022 [1]). However, not every Go package offers these benefits.One in particular exists which has named the Unsafe Package by Go. 
Go supports memory accesses and pointer operations exclusively through the Unsafe Package. 
The official Go website states that compatibility and protection by Go 1 compatibility guidelines may not be guaranteed by using this package [2]. 
Moreover, improper usage of it can lead to program crashes and unexpected behavior in programs due to misinterpretation of pointers (Costa et al., 2021 [3])<br>
The powerful type system in Go ensures stability, but it has been limited in some ways. 
Due to this, users commonly convert types by using Unsafe Package, serializing data, or modifying the value or obtaining the value of slice arrays with pointers [3]. 
There have been a few papers that recommend Generics be used in these types of situations ([3], Launiger et al., 2020 [4]). 
In March 2022, the Go development team announced the introduction of Generics through release 1.18. 
This has led to the availability of generic functions and structures within Go.
Related to the above, the datasets presented in the paper [4] were used to verify the refactoring of unsafe packages to generics. 
In their research, Launiger et al. stated that it would be possible to replace the use of the Unsafe package in Gorgonia with Generics and they selected codes for their dataset [4]. 
This paper investigates this aspect and examines the dataset selected by the authors.The dataset, which in this paper examined, are based on an open source project, which is named Gorgonia. 
A large percentage of the open-source project Gorgonia has been developed in Go (Go 96.4%, C 2.8%, other 0.8%). This project involves the realization of multi-dimensional arrays through serialization. 
There are several characteristics of Generics visible in this project. 
For examination, I filtered the dataset with the recommendation to use Generics by Go developer teams to check whether the usages of Unsafe Package in the dataset can be replaced by Generics.
Then, the dataset was analyzed structurally and semantically for this purpose. 
As part of the structural analysis, I converted the dataset into an AST tree structure using the Go ast package. 
In order to conduct the semantic analysis, I examined the comments and test codes in the project and manually evaluated the code. 
In order to eliminate the use of Unsafe Package, the patterns found are converted into Generics, and generic samples are then created according to the pattern. 
The patterns and generic samples will form the basis for refactoring. 
In order to address the above points, a tool has been developed. 
The tool takes as input a single Go file which is first converted to an AST structure. 
The code is then checked to see if it contains reused codes for replacing Generics. 
In cases where the codes have reused parts that match the patterns above, it will be examined to see if there is a use of Unsafe Packages. 
If that’s the case, the tool replaces the code used in the unsafe package with generics.