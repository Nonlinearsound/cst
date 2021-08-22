# cst
**cst - the command shell template parser**

cst is a command line template parser that uses comma seperated files (csv) as it's data source.
It is being used to transform data, present in csv files, into structured text files, such as HTML, JSON, XML or yaml files or any other structured file format.

Instead of using several tools like awk or writing a script for such a transformation every time you need it, it might be handy if you got a tool like cst at your hands.
Using cst, you just define a template file, choose certain comma seperated text files as your data and combined both to produce a structured text file that contains that data.

By this you can create XML or HTML files for your web services or web sites, configuraton files for all sorts of software, based on csv files that you either wrote byhand or exported, using some sort of software.

Usage:

cst -i \<input filepath\> -o \<output filepath\> 

The program needs an input and an output file path.
All needed data file paths are being defined in the template file itself in so called block definitions
The input file defines the template for the output file.
All parsed template definitions are processed and written to the output file.
