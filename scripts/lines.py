#!/usr/bin/python3

import re
import os
import subprocess
import argparse

PWD = os.path.dirname(os.path.abspath(__file__))
ROOT_DIR = os.path.dirname(PWD)
SRC_DIR = ROOT_DIR + "/src"

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

def count_lines(folder, pattern, verbose=True):
    print("Folder: {}, pattern: {}".format(folder, pattern))
    files = recurse_dir(folder)
    a = [re.fullmatch(pattern, f) for f in files]
    c = [l.string for l in filter(lambda a: a, a)]
    c.sort()

    cmd = ["wc"]
    cmd.extend(c)
    p = subprocess.Popen(cmd, stdout=subprocess.PIPE)
    out, err = p.communicate()

    out = out.decode("utf-8")
    if verbose:
        print(out)
    else:
        print(out.split("\n")[-2])

def list_all(verbose=False):
    # HTML
    print("HTML")
    count_lines(ROOT_DIR + "/client", ".*\.html$", verbose)

    # JS
    print("\nJS [lib, client, server]")
    count_lines(ROOT_DIR + "/lib", ".+\\.js$")
    count_lines(ROOT_DIR + "/server", ".+\\.js$", verbose)
    count_lines(ROOT_DIR + "/client", ".+\\.js$", verbose)

    # JS Tests
    print("\nJS Tests")
    count_lines(ROOT_DIR + "/tests", ".+\\.js", verbose)


    # Selenium tests
    # print("\nSelenium tests (Python)")
    # count_lines(ROOT_DIR + "tests", ".+\\.py$", verbose)

    # Misc, Wiki etc.
    print("\nWiki")
    count_lines(ROOT_DIR + "/wiki", ".+\\.md$", verbose)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="List the files matching a certain pattern and the number of contained lines, words and characters.")
    parser.add_argument("-p", "--pattern", help="pattern to use when looking for files")
    parser.add_argument("-d", "--dir", help="directory to look for files in")
    parser.add_argument("-v", "--verbose", help="lists all matched files", action='store_true')

    args = parser.parse_args()
    if not args.pattern:
        list_all(verbose=args.verbose)
    elif args.pattern and args.dir:
        count_lines(args.dir, args.pattern)
    else:
        print("Missing --pattern or --dir")
