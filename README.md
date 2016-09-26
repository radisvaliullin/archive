# test_task_5
The tesk task 5.

# Requirements:
python 3.5.x

# To Run:
python3 tasks_list.txt (or your file)

# Comment for output:
If map have size bigger than 5 x 5, script run more than 5 secons.

# Description:
Task 1. Next Biggest Number
Your test case is 1 number: N
Print out the next number greater than N using the same digits. For example, N=36 returns 63. If no permutation is greater (e.g. for N=77) print 0.

Task 2. Crawling Robot
Your test case is 6 numbers: N, M, X0, Y0, X1, and Y1 (in this order).
There’s a crawling robot inside a matrix of (N rows) × (M columns). The coordinates X [Y] denote zero-based index of a cell within a row [column]. The robot is placed at a start point with coordinates (X0, Y0) and needs to get to a finish point (X1, Y1). It can only move up, down, left or right between adjacent matrix cells and can visit the same cell only once. Robot cannot cross the borders of the matrix. Print out the number of all possible ways in which the robot can travel from the point of start to the point of finish (the output is one number).

Task 3. String Patterns
Your test case is 1 string: S
Locate all patterns in a string S (consisting of Arabic numerals and lowercase Latin letters `a` through `z` only). A pattern is defined as 2 or more adjacent characters within a string repeating at least twice. Print out a list of all existing patterns within a string in the order of their length, or `None` if none were found. For example, the output for string `a1a1a` would be `a1a`, `a1`, `1a`.

Notes:
• Input is provided in a text file with one or more lines, each line listing data for one test case (numbers separated with spaces, unless otherwise noted).
• Input filename supplied as (the only) command-line argument from the shell.
• Output: the program should produce output for every test case on a separate line in the shell, followed by a line of dashes `-----`.
• Your submissions are executed in a *NIX type environment. Assume softwares/shells etc are in their standard locations. Nothing else.
• Assume using Python 3.5 with built-in libraries only.
• Your program will be killed after 5 seconds of runtime.
• You don’t have to make the code fool-proof; assume all test cases are reasonable.
• You can use resources at your disposal (Google, etc.), except direct help from other people.
