import os, subprocess
from genericpath import isdir, isfile

# Write library imports to include.go file in /lib. Since Go plugins are
# only available on linux this is the next best option. Library dependencies
# are automatically added with the name of the folder. (this ofcourse means
# that the package name should be the same). Otherwise an undefined variable
# error will be raised.

# LIBRARY FORMATTING:
# Library files only need to have a 'Include' variable public (map of string
# and interface, where the string is the name of the exported function and
# the interface is the function itself).

filename = "lib/include.go"

if __name__ == "__main__":
    # Create file if it doesnt exist
    if not isfile(filename):
        os.open(filename, os.O_CREAT)
    
    with open(filename, "w+") as f:
        f.write("package lib\n")
        f.write("func init() {\n")

        contents = os.listdir("lib")
        for dirname in contents:
            if not isdir(f"lib/{dirname}") or dirname == "_libdump":
                continue
                
            # Go will raise error at compile time if there is something wrong here
            f.write(f'Add("{dirname}", {dirname}.Includes)\n')
        
        f.write("}")
        f.close()

    # Finally add the imports and format. Use subprocess to wait for finish
    subprocess.call("goimports -w lib")
    subprocess.call("gofmt -w -s -d lib")