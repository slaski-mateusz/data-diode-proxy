from sys import argv
from argparse import ArgumentParser
from yaml import safe_load, YAMLError
from schema import Schema, SchemaError, Optional
from time import sleep, time
from random import randint
from os import path
from os import scandir as sd

FLEN = 8
NUMLEN = 23

CONF_SCHEMA = Schema({
    "out-file": {
        "path": str,
        "name": str,
        "ext": str,
        "min-size-kb": int,
        "max-size-kb": int
    },
    "cycle": {
        "seconds": int,
        "random-sec-offset": int
    },
    Optional("errors"): {
        "skip-file-chance": int,
        "skip-line-chance": int
    }
})


def to_skip(chance):
    probe = randint(0, 100)
    if probe < chance:
        return True
    else:
        return False


def main():
    default_conf_name = argv[0].replace("py", "yaml")
    a_parser = ArgumentParser()
    a_parser.add_argument("-c", "--conf", type=str, default=default_conf_name)
    a_parsed = a_parser.parse_args()
    conf_file_name = a_parsed.conf
    skip_file_chance = 0
    skip_line_chance = 0
    with open(conf_file_name, "r") as conf_file:
        try:
            conf_data = safe_load(conf_file)
        except YAMLError as yaml_exception:
            raise yaml_exception
    try:
        CONF_SCHEMA.validate(conf_data)
    except SchemaError as schema_exception:
        raise schema_exception
    out_path = conf_data["out-file"]["path"]
    out_name = conf_data["out-file"]["name"]
    out_ext = conf_data["out-file"]["ext"]
    out_min_lines = conf_data["out-file"]["min-size-kb"] * 32
    out_max_lines = conf_data["out-file"]["max-size-kb"] * 32
    if "errors" in conf_data.keys():
        skip_file_chance = conf_data["errors"]["skip-file-chance"]
        skip_line_chance = conf_data["errors"]["skip-line-chance"]
    file_count = 0
    existing_files = sd(out_path)
    for file_to_check in existing_files:
        filename = file_to_check.name.split(".")[0]
        filenumtxt = filename.removeprefix(conf_data["out-file"]["name"])
        try:
            filenum = int(filenumtxt)
            if filenum > file_count:
                file_count = filenum + 1
        except ValueError:
            continue
    while True:
        out_filename = out_name + str(file_count).zfill(FLEN) + "." + out_ext
        if to_skip(skip_file_chance):
            print("File {} skipped".format(out_filename))
        else:
            out_filepath = path.join(
                out_path,
                out_filename
                )
            out_lines = randint(
                out_min_lines,
                out_max_lines
                )
            with open(out_filepath, "w") as out_file:
                line_count = 0
                written_line_count = 0
                while written_line_count < out_lines:
                    if to_skip(skip_line_chance):
                        pass
                    else:
                        out_file.write(
                            str(file_count).zfill(FLEN) + str(line_count).zfill(NUMLEN) + "\n"
                            )
                        written_line_count += 1
                    line_count += 1
            print("file {} written".format(out_filename))
        file_count += 1
        rand_sleep_offs = randint(
            0,
            conf_data["cycle"]["random-sec-offset"]
            )
        sleep(conf_data["cycle"]["seconds"] + rand_sleep_offs)


if __name__ == "__main__":
    main()
