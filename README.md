# cst
**cst - the command shell template parser**

cst is a command line template parser that uses comma seperated files (csv) as it's data source.
It is being used to transform data, present in csv files, into structured text files, such as HTML, JSON, XML or yaml files or any other structured file format.

Instead of using several tools like awk or writing a script for such a transformation every time you need it, it might be handy if you got a tool like cst at your hands.
Using cst, you just define a template file, choose certain comma seperated text files as your data and combined both to produce a structured text file that contains that data.

By this you can create XML or HTML files for your web services or web sites, configuraton files for all sorts of software, based on csv files that you either wrote byhand or exported, using some sort of software.

## Usage
-----
cst -i \<input filepath\> -o \<output filepath\> 

The program needs an input and an output file path.
All needed data file paths are being defined in the template file itself in so called block definitions
The input file defines the template for the output file.
All parsed template definitions are processed and written to the output file.

## Example

### Template blocks

This is an example of a template file and to be passed as an input file with the argument -i.
Notice the definitions in brackets, like "store", "block-start" and "block-end". These are template block definitions that the parser and the template engine will use to replace them with content of the specified data csv file.

Also there are template placeholder definitions like {{name}} specified. These are simple template definitions where the token is replaced with the value of the specified key.

For simplicity, the key-value store and the data files for the template blocks are simple CSV files without headlines like you might have exported from applications like Excel or LibreOffice.
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

In the following example we want to use the template engine to create an HTML file that uses two csv files to create two different unsorted list definitions in the HTML source.
For that, we define two "foreach" blocks in the form of
```
(block-start;type:foreach;source:\<datafile-path\>)
 templated string with placeholder tokens for the column of the data file
 for instance {0} for the first column. Use as many columns you have in the data file.
(block-end)
```
"block-end" is important to tell the templating parser that the block ends here.
The template engine will then read in the specified data file and loop over its lines/rows. For every row it will loop over as many template lines are being defined in the template block and replace the placeholders with the content of the respective column of each line/row in the data file.

If we have a data.csv file with the content
```
foo,bar,name,fred
baz,test,name,walter
```
and a template block definition like
```
(block-start;type:foreach;source:data.csv)
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

### Key-Value Templates

In every line of your tempalte file you can define placeholders for keys that are defined in the key-value csv file, you specified as following:
```
(store;source:<key-value-file-path>)
```
Now, you can use that data by defining a placeholder like this
```
Some text line with the value of the key foo = {{foo}}.
Another line where we will put the value of {{bar}} into the string.
```
Then, if we got this csv file as our key-value data file
```
foo,fred
bar,walter
```
we will end up with these two lines in our output file
```
Some text line with the value of the key foo = fred.
Another line where we will put the value of walter into the string.
```
You can use these placeholders everywhere in the source file. They can even be used inside block definitions. Have a look in the Example folder in the source.txt file.

### JSON export and import

JSON is a structured data format that is being supported by almost all software that deals with transforming and/or transferring data. If you want to use JSON for defining the blocks, you can use the parser to parse the course file and export the block structure as JSON. You can then use that structure to create ana adaptation of that definition or create your very own definition for the templating engine to be used as input.

You can export the parsed block definitions as JSON to a file using the argument -outputjson. The JSON data is written to ./blocksDefinitions.json by default. If you want to write to a different path, use the -jsonoutfile argument. Once done or written by hand, you can use such JSON file as the input source for your block definitions using the -jsoninput argument. 

cst -i ./source.txt -outputjson -jsonoutfile=./json-definitions.json

will create the JSON file json-definitions.json in the same directory as cst itself resides. 

cst -inputjson -jsoninput=./json-definitions.json

then reads that JSON file and creates the block definitions out of that and uses it for the templating engine.
