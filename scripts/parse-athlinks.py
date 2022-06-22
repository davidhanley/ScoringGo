import sys
import re

inf = sys.argv[1]

lines = [ l.strip() for l in open(inf,"r") ]

index = 0 
while index < len(lines):
  line = lines[index]
  line = line.split()
  if len(line)>=3 and "Bib" in line[1]:
    name = lines[index-1]
    gender = line[0]
    age = line[1].split("Bib")[0]
  
    parts = ["",name,age,gender]
    print ",".join(parts)
  index = index + 1 

   
