# Bug Reproduction

## 包的性质

当前 test_model_fix 保存的是被测模型修复后的结果源码，不是初始含 Bug 源码。要复现原始缺陷，必须检出下面固定的 parent SHA，再应用仓库内的可信测试补丁；不要在当前修复结果源码上期待重新出现修复前失败。

## 问题现象

帮我调查 CSV SQL 的 RIGHT JOIN 异常，先不要修改代码。

左右表都有重复键和各自独有键时，匹配键的多对多结果正确，但结果保留的是左表独有行并给右列补 NULL；右表独有行反而消失。WHERE 过滤右表独有行得到空结果，按右表键 GROUP BY 后也出现 NULL 分组和错误计数。

请给出具体 Go 文件、具体符号、错误行为及其造成这些症状的完整因果链，并提供实际代码或复现证据。先只提交诊断结论，不要改仓库代码。

## 含 Bug 版本

- 仓库：zhanglei10281852-gif/gogo-07
- 仓库地址：https://github.com/zhanglei10281852-gif/gogo-07.git
- parent SHA：8f8498a16011a0dcc863e4085ced62a7e7b94689

## 复现步骤

```bash
git clone -- https://github.com/zhanglei10281852-gif/gogo-07.git bug-repro
cd bug-repro
git checkout --detach 8f8498a16011a0dcc863e4085ced62a7e7b94689
git apply ../BENZHI_VALIDATION/trusted-test.patch
go test ./internal/executor -run "^TestRightJoin" -count=1 -v
```

## 双架构完整错误信息

### linux/amd64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/executor -run "^TestRightJoin" -count=1 -v
=== RUN   TestRightJoinPreservesRightRowsAndNullExtendsLeft
    right_join_test.go:77: 
        	Error Trace:	/app/internal/executor/right_join_test.go:77
        	Error:      	Not equal: 
        	            	expected: [][]string{[]string{"1", "L1a", "1", "R1a"}, []string{"1", "L1a", "1", "R1b"}, []string{"1", "L1b", "1", "R1a"}, []string{"1", "L1b", "1", "R1b"}, []string{"NULL", "NULL", "3", "R3"}, []string{"NULL", "NULL", "5", "R5"}}
        	            	actual  : [][]string{[]string{"2", "L2", "NULL", "NULL"}, []string{"4", "L4", "NULL", "NULL"}, []string{"1", "L1a", "1", "R1a"}, []string{"1", "L1a", "1", "R1b"}, []string{"1", "L1b", "1", "R1a"}, []string{"1", "L1b", "1", "R1b"}}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1,2 +1,14 @@
        	            	 ([][]string) (len=6) {
        	            	+ ([]string) (len=4) {
        	            	+  (string) (len=1) "2",
        	            	+  (string) (len=2) "L2",
        	            	+  (string) (len=4) "NULL",
        	            	+  (string) (len=4) "NULL"
        	            	+ },
        	            	+ ([]string) (len=4) {
        	            	+  (string) (len=1) "4",
        	            	+  (string) (len=2) "L4",
        	            	+  (string) (len=4) "NULL",
        	            	+  (string) (len=4) "NULL"
        	            	+ },
        	            	  ([]string) (len=4) {
        	            	@@ -24,14 +36,2 @@
        	            	   (string) (len=3) "R1b"
        	            	- },
        	            	- ([]string) (len=4) {
        	            	-  (string) (len=4) "NULL",
        	            	-  (string) (len=4) "NULL",
        	            	-  (string) (len=1) "3",
        	            	-  (string) (len=2) "R3"
        	            	- },
        	            	- ([]string) (len=4) {
        	            	-  (string) (len=4) "NULL",
        	            	-  (string) (len=4) "NULL",
        	            	-  (string) (len=1) "5",
        	            	-  (string) (len=2) "R5"
        	            	  }
        	Test:       	TestRightJoinPreservesRightRowsAndNullExtendsLeft
    right_join_test.go:86: 
        	Error Trace:	/app/internal/executor/right_join_test.go:86
        	Error:      	Should be true
        	Test:       	TestRightJoinPreservesRightRowsAndNullExtendsLeft
        	Messages:   	unmatched right rows must NULL-extend left columns
    right_join_test.go:87: 
        	Error Trace:	/app/internal/executor/right_join_test.go:87
        	Error:      	Should be true
        	Test:       	TestRightJoinPreservesRightRowsAndNullExtendsLeft
        	Messages:   	unmatched right rows must NULL-extend left columns
    right_join_test.go:86: 
        	Error Trace:	/app/internal/executor/right_join_test.go:86
        	Error:      	Should be true
        	Test:       	TestRightJoinPreservesRightRowsAndNullExtendsLeft
        	Messages:   	unmatched right rows must NULL-extend left columns
    right_join_test.go:87: 
        	Error Trace:	/app/internal/executor/right_join_test.go:87
        	Error:      	Should be true
        	Test:       	TestRightJoinPreservesRightRowsAndNullExtendsLeft
        	Messages:   	unmatched right rows must NULL-extend left columns
--- FAIL: TestRightJoinPreservesRightRowsAndNullExtendsLeft (0.00s)
=== RUN   TestRightJoinManyToManyMatches
--- PASS: TestRightJoinManyToManyMatches (0.00s)
=== RUN   TestRightJoinWhereSeesUnmatchedRightRows
    right_join_test.go:116: 
        	Error Trace:	/app/internal/executor/right_join_test.go:116
        	Error:      	Not equal: 
        	            	expected: [][]string{[]string{"3", "R3"}, []string{"5", "R5"}}
        	            	actual  : [][]string{}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1,10 +1,2 @@
        	            	-([][]string) (len=2) {
        	            	- ([]string) (len=2) {
        	            	-  (string) (len=1) "3",
        	            	-  (string) (len=2) "R3"
        	            	- },
        	            	- ([]string) (len=2) {
        	            	-  (string) (len=1) "5",
        	            	-  (string) (len=2) "R5"
        	            	- }
        	            	+([][]string) {
        	            	 }
        	Test:       	TestRightJoinWhereSeesUnmatchedRightRows
--- FAIL: TestRightJoinWhereSeesUnmatchedRightRows (0.00s)
=== RUN   TestRightJoinAggregationIncludesUnmatchedRightRows
    right_join_test.go:128: 
        	Error Trace:	/app/internal/executor/right_join_test.go:128
        	Error:      	Not equal: 
        	            	expected: [][]string{[]string{"1", "4"}, []string{"3", "0"}, []string{"5", "0"}}
        	            	actual  : [][]string{[]string{"NULL", "2"}, []string{"1", "4"}}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1,2 +1,6 @@
        	            	-([][]string) (len=3) {
        	            	+([][]string) (len=2) {
        	            	+ ([]string) (len=2) {
        	            	+  (string) (len=4) "NULL",
        	            	+  (string) (len=1) "2"
        	            	+ },
        	            	  ([]string) (len=2) {
        	            	@@ -4,10 +8,2 @@
        	            	   (string) (len=1) "4"
        	            	- },
        	            	- ([]string) (len=2) {
        	            	-  (string) (len=1) "3",
        	            	-  (string) (len=1) "0"
        	            	- },
        	            	- ([]string) (len=2) {
        	            	-  (string) (len=1) "5",
        	            	-  (string) (len=1) "0"
        	            	  }
        	Test:       	TestRightJoinAggregationIncludesUnmatchedRightRows
--- FAIL: TestRightJoinAggregationIncludesUnmatchedRightRows (0.00s)
FAIL
FAIL	csvsql/internal/executor	0.005s
FAIL

```

stderr：

```text
(empty)
```

### linux/arm64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/executor -run "^TestRightJoin" -count=1 -v
=== RUN   TestRightJoinPreservesRightRowsAndNullExtendsLeft
    right_join_test.go:77: 
        	Error Trace:	/app/internal/executor/right_join_test.go:77
        	Error:      	Not equal: 
        	            	expected: [][]string{[]string{"1", "L1a", "1", "R1a"}, []string{"1", "L1a", "1", "R1b"}, []string{"1", "L1b", "1", "R1a"}, []string{"1", "L1b", "1", "R1b"}, []string{"NULL", "NULL", "3", "R3"}, []string{"NULL", "NULL", "5", "R5"}}
        	            	actual  : [][]string{[]string{"2", "L2", "NULL", "NULL"}, []string{"4", "L4", "NULL", "NULL"}, []string{"1", "L1a", "1", "R1a"}, []string{"1", "L1a", "1", "R1b"}, []string{"1", "L1b", "1", "R1a"}, []string{"1", "L1b", "1", "R1b"}}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1,2 +1,14 @@
        	            	 ([][]string) (len=6) {
        	            	+ ([]string) (len=4) {
        	            	+  (string) (len=1) "2",
        	            	+  (string) (len=2) "L2",
        	            	+  (string) (len=4) "NULL",
        	            	+  (string) (len=4) "NULL"
        	            	+ },
        	            	+ ([]string) (len=4) {
        	            	+  (string) (len=1) "4",
        	            	+  (string) (len=2) "L4",
        	            	+  (string) (len=4) "NULL",
        	            	+  (string) (len=4) "NULL"
        	            	+ },
        	            	  ([]string) (len=4) {
        	            	@@ -24,14 +36,2 @@
        	            	   (string) (len=3) "R1b"
        	            	- },
        	            	- ([]string) (len=4) {
        	            	-  (string) (len=4) "NULL",
        	            	-  (string) (len=4) "NULL",
        	            	-  (string) (len=1) "3",
        	            	-  (string) (len=2) "R3"
        	            	- },
        	            	- ([]string) (len=4) {
        	            	-  (string) (len=4) "NULL",
        	            	-  (string) (len=4) "NULL",
        	            	-  (string) (len=1) "5",
        	            	-  (string) (len=2) "R5"
        	            	  }
        	Test:       	TestRightJoinPreservesRightRowsAndNullExtendsLeft
    right_join_test.go:86: 
        	Error Trace:	/app/internal/executor/right_join_test.go:86
        	Error:      	Should be true
        	Test:       	TestRightJoinPreservesRightRowsAndNullExtendsLeft
        	Messages:   	unmatched right rows must NULL-extend left columns
    right_join_test.go:87: 
        	Error Trace:	/app/internal/executor/right_join_test.go:87
        	Error:      	Should be true
        	Test:       	TestRightJoinPreservesRightRowsAndNullExtendsLeft
        	Messages:   	unmatched right rows must NULL-extend left columns
    right_join_test.go:86: 
        	Error Trace:	/app/internal/executor/right_join_test.go:86
        	Error:      	Should be true
        	Test:       	TestRightJoinPreservesRightRowsAndNullExtendsLeft
        	Messages:   	unmatched right rows must NULL-extend left columns
    right_join_test.go:87: 
        	Error Trace:	/app/internal/executor/right_join_test.go:87
        	Error:      	Should be true
        	Test:       	TestRightJoinPreservesRightRowsAndNullExtendsLeft
        	Messages:   	unmatched right rows must NULL-extend left columns
--- FAIL: TestRightJoinPreservesRightRowsAndNullExtendsLeft (0.08s)
=== RUN   TestRightJoinManyToManyMatches
--- PASS: TestRightJoinManyToManyMatches (0.00s)
=== RUN   TestRightJoinWhereSeesUnmatchedRightRows
    right_join_test.go:116: 
        	Error Trace:	/app/internal/executor/right_join_test.go:116
        	Error:      	Not equal: 
        	            	expected: [][]string{[]string{"3", "R3"}, []string{"5", "R5"}}
        	            	actual  : [][]string{}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1,10 +1,2 @@
        	            	-([][]string) (len=2) {
        	            	- ([]string) (len=2) {
        	            	-  (string) (len=1) "3",
        	            	-  (string) (len=2) "R3"
        	            	- },
        	            	- ([]string) (len=2) {
        	            	-  (string) (len=1) "5",
        	            	-  (string) (len=2) "R5"
        	            	- }
        	            	+([][]string) {
        	            	 }
        	Test:       	TestRightJoinWhereSeesUnmatchedRightRows
--- FAIL: TestRightJoinWhereSeesUnmatchedRightRows (0.00s)
=== RUN   TestRightJoinAggregationIncludesUnmatchedRightRows
    right_join_test.go:128: 
        	Error Trace:	/app/internal/executor/right_join_test.go:128
        	Error:      	Not equal: 
        	            	expected: [][]string{[]string{"1", "4"}, []string{"3", "0"}, []string{"5", "0"}}
        	            	actual  : [][]string{[]string{"NULL", "2"}, []string{"1", "4"}}
        	            	
        	            	Diff:
        	            	--- Expected
        	            	+++ Actual
        	            	@@ -1,2 +1,6 @@
        	            	-([][]string) (len=3) {
        	            	+([][]string) (len=2) {
        	            	+ ([]string) (len=2) {
        	            	+  (string) (len=4) "NULL",
        	            	+  (string) (len=1) "2"
        	            	+ },
        	            	  ([]string) (len=2) {
        	            	@@ -4,10 +8,2 @@
        	            	   (string) (len=1) "4"
        	            	- },
        	            	- ([]string) (len=2) {
        	            	-  (string) (len=1) "3",
        	            	-  (string) (len=1) "0"
        	            	- },
        	            	- ([]string) (len=2) {
        	            	-  (string) (len=1) "5",
        	            	-  (string) (len=1) "0"
        	            	  }
        	Test:       	TestRightJoinAggregationIncludesUnmatchedRightRows
--- FAIL: TestRightJoinAggregationIncludesUnmatchedRightRows (0.02s)
FAIL
FAIL	csvsql/internal/executor	0.529s
FAIL

```

stderr：

```text
(empty)
```

## 通过条件

准确定位 internal/executor/executor.go 的 (*Executor).materializeJoin
解释 RIGHT JOIN 被按左侧保留语义执行、未追踪和补出未匹配右行的完整机制
结论有代码阅读或定向复现证据，目标仓库保持零改动
