import os

# Generates docs for libraries from the function docsstrings.
# Todo: make this not scuffed as fuck (also comment)

def main():
    c = os.listdir("./lib")
    for foldername in c:
        if len(foldername.split(".")) != 1:
            continue
        
        path = f"./lib/{foldername}"
        files = os.listdir(path)
        total = []
        for file in files:
            if not file.endswith(".go"):
                continue
        
            content = open(f"{path}/{file}").readlines()
            total.extend(add_doc(content))

        if len(total) != 0:
            create_file(path, foldername, total)


def add_doc(c: str):
    docs = []
    start = 0
    for idx, l in enumerate(c):
        notab = l.replace("\t", "")
        if notab.startswith("/*") and l.count("*/") == 1:
            docs.append(("", f"{notab[3:len(notab)-3]}\n"))
            continue

        if l.startswith("/*"):
            start = idx
        if l.startswith("*/"):
            interval = c[start+1:idx]
            doc = "".join(interval[:len(interval)-1]).replace("\t", "")
            func = interval[len(interval)-1].replace("\t", "")
            docs.append((doc, func))
    
    return docs


def create_file(path: str, libname: str, docs: list):
    filename = f"{path}/{libname}_docs.md"
    os.open(filename, os.O_CREAT)
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