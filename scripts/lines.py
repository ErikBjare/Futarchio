#!/usr/bin/python3

import re
import os
import subprocess
import argparse
def recurse_dir(folder):
    files = [folder + "/" + f for f in os.listdir(folder)]
    dirs = list(filter(lambda x: os.path.isdir(x), files))
    files = list(filter(lambda x: os.path.isfile(x), files))
    for dir in dirs:
        new_files = recurse_dir(dir)
        if len(new_files) == 0:
            print("Found empty directory {}".format(dir))
        files.extend(new_files)
    return files

site_src = recurse_dir("src/site/src")

def count_lines(folder, pattern):
    print("Folder: {}, pattern: {}".format(folder, pattern))
    files = recurse_dir(folder)
    a = [re.fullmatch(pattern, f) for f in files]
    c = [l.string for l in filter(lambda a: a, a)]
    c.sort()
    
    cmd = ["wc"]
    cmd.extend(c)
    subprocess.call(cmd)

def list_all():
    # Site, HTML & JS
    print("HTML")
    count_lines("src/site/src", ".*\.html$")
    count_lines("src/site/src", ".*[^(\.min)]\\.js$")

    # Go, backend and tests
    print("\nGo")
    count_lines("src", ".*[^(_test)]\.go$")
    print("\nGo tests")
    count_lines("src", ".*_test\.go$")

    # Misc, Wiki etc.
    print("\nWiki")
    count_lines("wiki", ".*\\.md$")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="List the files matching a certain pattern and the number of contained lines, words and characters.")
    parser.add_argument("-p", "--pattern", help="pattern to use when looking for files")
    parser.add_argument("-d", "--dir", help="directory to look for files in")

    args = parser.parse_args()
    if not args.pattern:
        list_all()
    elif args.pattern and args.dir:
        count_lines(args.dir, args.pattern)
    else:
        print("Missing --pattern or --dir")
