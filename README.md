# csvmung

csvmung is a simple tool to manipulate csv files. It can perform the following transformations:

* Filter out rows that don't match a regular expression.
* Output only rows with unique values.
* Output only certain columns.
* Split columns based on a regular expression.

## Transforms

You can use the following transformations. All column references are zero indexed.

### re

	re:<column>:<pattern>
	
Only pass through rows where the `column` matches the `pattern`.

### split

	split:<column>:<pattern>
	
Split `column` into separate columns using `pattern`.


### uniq

	uniq:<column>
	
Only pass through rows where the value in `column` hasn't yet been seen.

### cols

	cols:<col1>:<col2>:...:<colN>

Only pass through the specified columns. If the column specifier isn't an integer, the value will be passed through in that position.

## Examples

Given a csv file with the following contents, named foo.csv:

    a,b-j,c,d
    a,k-n,l,m
    a,n-f,g,x
    a,p-p,d,d

You can get all the rows that contain 'd' in the fourth column:

	$ csvmung -i foo.csv re:3:'^d$'
	a,b-j,c,d
    a,p-p,d,d

or only output the columns with a unique value in the fourth column:
    
 	$ csvmung -i foo.csv uniq:3
    a,b-j,c,d
    a,k-n,l,m
    a,n-f,g,x

or output just the first and fourth columns:

 	$ csvmung -i foo.csv cols:0:3
    a,d
    a,m
    a,x
    a,d

or split the value in the second column:

 	$ csvmung -i foo.csv split:1:'-'
    a,b,j,c,d
    a,k,n,l,m
    a,n,f,g,x
    a,p,p,d,d

Finally, you can chain the operators so that each works on the output of the last:

 	$ csvmung -i foo.csv split:1:'-' re:4:'^d$'
    a,b,j,c,d
    a,p,p,d,d

