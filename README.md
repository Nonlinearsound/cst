# cst
**cst - the command shell template parser**

cst is a command line template parser that uses comma seperated files (csv) as it's data source.
It is being used to transform data, present in csv files, into structured text files, such as HTML, JSON, XML or yaml files or any other structured file format.

Instead of using several tools like awk or writing a script for such a transformation every time you need it, it might be handy if you got a tool like cst at your hands.
Using cst, you just define a template file, choose certain comma seperated text files as your data and combined both to produce a structured text file that contains that data.

By this you can create XML or HTML files for your web services or web sites, configuraton files for all sorts of software, based on csv files that you either wrote byhand or exported, using some sort of software.

Usage
-----
cst -i \<input filepath\> -o \<output filepath\> 

The program needs an input and an output file path.
All needed data file paths are being defined in the template file itself in so called block definitions
The input file defines the template for the output file.
All parsed template definitions are processed and written to the output file.

Example
-----

Create a tempalte file and pass it as an input file with the argument -i

```
(store;source:Example/keystore.csv)
<body>
<ul id="list-{{name}}">
(block-start;type:foreach;source:Example/data.csv)
 <li class="{{classname}}"><p>{0}</p><p>{1}<p/></li>
 <h1>{0}</h1>
(block-end)
</ul>

<ul id="list-{{name2}}">
(block-start;type:foreach;source:Example/data2.csv)
 <li class="{{classname}}"><p>{0}</p><p>{1}<p/><span>{2}<span/></li>
(block-end)
</ul>
</body>
```

Here we want to use the template engine to create an HTML file that uses two csv files to create two different unsorted list definitions in the HTML source.
For that, we define two "foreach" blocks in the form of
```
(block-start;type:foreach;source:\<datafile-path\>)
 templated string with placeholder tokens for the column of the data file
 for instance {0} for the first column. Use as many columns you have in the data file.
(block-end)
```
"block-end" is importand to tell the templating parser that the block ends here.
The template engine will then read in the specified data file and loop over its lines/rows. For every row it will loop over as many template lines are being defined in the template block and replace the placeholders with the content of the respective column of each line/row in the data file.

If we have a data.csv file with the content
```
foo,bar,name,fred
baz,test,name,walter
```
and a template block definition like
```
(block-start;type:foreach;source:Example/data.csv)
 <h2 id="{0}">{1}</h2>
 <p id="{2}">{3}</p>
(block-end)
```
we will end up with an output file containing the lines
```
 <h2 id="foo">bar</h2>
 <p id="name">fred</p>
 <h2 id="baz">test</h2>
 <p id="name">walter</p>
```
You can create as many block as you wish and specifiy different data files for them so that you can out as many data as you like in the structured output text file.

