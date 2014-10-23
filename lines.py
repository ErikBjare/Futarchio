#!/usr/bin/python3

import re
import os
import subprocess

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
