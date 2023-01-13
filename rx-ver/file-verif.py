from sys import argv
from os import scandir as sd
from os.path import join as join_path
from yaml import safe_load, YAMLError
from schema import Schema, SchemaError
from time import time, sleep

FLEN = 8
NUMLEN = 23

CONF_SCHEMA = Schema({
    "in-file": {
        "path": str,
        "name": str,
        "ext": str
    },
    "cycle": {
        "seconds": int
    }
})


def main():
    conf_file_name = argv[0].replace("py", "yaml")
    with open(conf_file_name, "r") as conf_file:
        try:
            conf_data = safe_load(conf_file)
        except YAMLError as yaml_exception:
            raise yaml_exception
    try:
        CONF_SCHEMA.validate(conf_data)
    except SchemaError as schema_exception:
        raise schema_exception
    while True:
        with sd(conf_data["in-file"]["path"]) as in_files:
            for file_to_check in in_files:
                if file_to_check.is_file():
                    file_mtime = file_to_check.stat().st_mtime
                    if file_mtime > time() - conf_data["cycle"]["seconds"]:
                        file_lines = []
                        file_to_check_path = join_path(
                            conf_data["in-file"]["path"],
                            file_to_check.name
                        )
                        with open(file_to_check_path) as file_to_read:
                            file_lines = file_to_read.readlines()
                        filename = file_to_check.name.split(".")[0]
                        filenumtxt = filename.removeprefix(conf_data["in-file"]["name"])
                        for idx, line in enumerate(file_lines):
                            fpart = line[0:8]
                            if fpart != filenumtxt:
                                print(
                                    "File {} line number contains value {} instead of {}".format(
                                        idx,
                                        fpart,
                                        filename
                                    )
                                )
                            npart = line[8:]
                            try:
                                num_in_line = int(npart)
                            except ValueError:
                                print(
                                    "Line {} not contains {} instead of number".format(
                                        idx,
                                        npart
                                    )
                                )
                            if idx + 1 < len(file_lines):
                                next_npart = file_lines[idx + 1][8:]
                                try:
                                    num_in_next_line = int(next_npart)
                                except ValueError:
                                    continue
                                if num_in_next_line != num_in_line + 1:
                                    print(
                                        "In file '{}' numbers in lines {} and {} are not ascending by 1. They are {} and {}".format(
                                            filename,
                                            idx,
                                            idx+1,
                                            num_in_line,
                                            num_in_next_line
                                        )
                                    )
        sleep(conf_data["cycle"]["seconds"])


if __name__ == "__main__":
    main()
