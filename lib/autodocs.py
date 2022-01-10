import os

# Generates docs for libraries from the function docstrings.

def main():
    libraries = os.listdir("./lib")
    for foldername in libraries:
        # Check for file not directory
        if len(foldername.split(".")) != 1:
            continue
        
        # Get library files in folder
        path = f"./lib/{foldername}"
        files = os.listdir(path)

        total = []
        for file in files:
            # Must be Go file
            if not file.endswith(".go"):
                continue
            
            # Read file and add formatted comments to total
            content = open(f"{path}/{file}").readlines()
            total.extend(add_doc(content))

        if len(total) != 0:
            # Write doc content to md file
            create_file(path, foldername, total)


# Returns formatted go comment as markdown section
def add_doc(c: str):
    # Total comments/messages
    docs = []

    start = 0
    for idx, l in enumerate(c):
        # Remove tabs
        raw = l.replace("\t", "")
        # One liner comment
        if raw.startswith("/*") and l.count("*/") == 1:
            # Append content bewteen pairs
            docs.append(("", f"{raw[3:len(raw)-3]}\n"))
            continue
        
        # Multiline comment, get content between pairs
        if l.startswith("/*"):
            start = idx
        if l.startswith("*/"):
            interval = c[start+1:idx]
            # Remove tabs
            doc = "".join(interval[:len(interval)-1]).replace("\t", "")
            func = interval[len(interval)-1].replace("\t", "")
            docs.append((doc, func))
    
    return docs


# Write docs to file (create file)
def create_file(path: str, libname: str, docs: list):
    # Create file if it doesnt exist already
    filename = f"{path}/{libname}_docs.md"
    os.open(filename, os.O_CREAT)

    # Open and write formatted with markdown
    with open(filename, "w+") as f:
        f.write(f"# Methods in {libname} library\n\n")
        for func in docs:
            name = func[1].split(" ")[1].split("(")[0]
            f.write(f"## **`{name}`**\n\n")
            f.write(f"{func[0]}\n")
            f.write(f"```go\n{func[1]}```\n\n")
            f.write("<br>\n\n")


if __name__ == "__main__":
    main()