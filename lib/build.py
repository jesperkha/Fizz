import os, subprocess
from genericpath import isdir, isfile

if __name__ == "__main__":
    # Create init functions for libraries. This is to export the public
    # functions through the Include map. Creates a new 'export.go' file
    # (temporary) where the functions are included automatically, where
    # only functions with capital first letters (go standard) are exported.
    libs = os.listdir("lib")
    for lib in libs:
        if not isdir(f"lib/{lib}"):
            continue
        
        for file in os.listdir(f"lib/{lib}"):
            # Directory or not go file
            if len(file.split(".")) == 1 or file.split(".")[1] != "go":
                continue

            exports: list[str] = [] # list of function names
            with open(f"lib/{lib}/{file}", "r+") as f:
                lines = f.readlines()
                for line in lines:
                    if not line.startswith("func"):
                        continue

                    # First char in function name
                    is_export = line.split(" ")[1][0].isupper()
                    if not is_export:
                        continue
                        
                    func_name = line.split(" ")[1].split("(")[0]
                    exports.append(func_name)
            
            if len(exports) == 0:
                continue

            export_file = f"lib/{lib}/export.go"
            if not isfile(export_file):
                os.open(export_file, os.O_CREAT)
            
            with open(export_file, "w+") as f:
                f.write("// AUTO-GENERATED FOR BUILD, DO NOT EDIT\n")
                f.write("// https://github.com/jesperkha/Fizz/blob/main/docs/libraries.md\n")
                f.write(f"package {lib}\n")
                f.write("var Includes = map[string]interface{}{}\n")
                f.write("func init() {\n")
                for name in exports:
                    lower_name = name[0].lower() + name[1:]
                    f.write(f'Includes["{lower_name}"] = {name}\n')
                f.write("}")
                
    # Write library imports to include.go file in /lib. Since Go plugins are
    # only available on linux this is the next best option. Library dependencies
    # are automatically added with the name of the folder. (this ofcourse means
    # that the package name should be the same). Otherwise an undefined variable
    # error will be raised.
    filename = "lib/include.go"
    if not isfile(filename):
        os.open(filename, os.O_CREAT)
    
    with open(filename, "w+") as f:
        f.write("package lib\n")
        f.write("func init() {\n")

        count = 0
        contents = os.listdir("lib")
        for dirname in contents:
            if not isdir(f"lib/{dirname}") or dirname == "_libdump":
                continue
                
            # Go will raise error at compile time if there is something wrong here
            if isfile(f"lib/{dirname}/export.go"):
                f.write(f'Add("{dirname}", {dirname}.Includes)\n')
                count += 1
            else:
                print(f"[WARNING] Library '{dirname}' doesnt export any functions")
        
        f.write("}")

    # Finally add the imports
    cmd_imports = "goimports -w lib"
    print(f"[CMD] {cmd_imports}")
    try:
        subprocess.call(cmd_imports)
    except:
        print("[FATAL] Failed to run goimports. Install:\ngo install golang.org/x/tools/cmd/goimports@latest")
        exit(1)
        
    print(f"[INFO] Included {count} libraries")