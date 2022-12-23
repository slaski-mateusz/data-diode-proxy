from sys import argv
from argparse import ArgumentParser
import yaml
from schema import Schema, SchemaError
from time import sleep, time
from random import randint
from os import path

conf_schema = Schema({
    "out-file": {
        "path": str,
        "name": str,
        "ext": str,
        "min-size-mb": int,
        "max-size-mb": int
    },
    "cycle": {
        "seconds": int,
        "random-sec-offset": int
    }
})


def main():
    a_parser = ArgumentParser()
    a_parser.add_argument("-c", "--conf", type=str)
    a_parsed = a_parser.parse_args()
    conf_file_name = argv[0].replace("py", "yaml")
    if a_parsed.conf is not None:
        conf_file_name = a_parsed.conf
    with open(conf_file_name, "r") as conf_file:
        try:
            conf_data = yaml.safe_load(conf_file)
        except yaml.YAMLError as yaml_exc:
            raise yaml_exc
    try:
        conf_schema.validate(conf_data)
    except SchemaError as sche_exc:
        raise sche_exc
    out_path = conf_data["out-file"]["path"]
    out_name = conf_data["out-file"]["name"]
    out_ext = conf_data["out-file"]["ext"]
    out_min_lines = conf_data["out-file"]["min-size-mb"] * 1024 * 32
    out_max_lines = conf_data["out-file"]["max-size-mb"] * 1024 * 32
    file_count = 0
    while True:
        out_filename = out_name + str(file_count).zfill(8) + "." + out_ext
        out_filepath = path.join(
            out_path,
            out_filename
            )
        out_lines = randint(
            out_min_lines,
            out_max_lines
            )
        with open(out_filepath, "w") as out_file:
            for i in range(out_lines):
                out_file.write(
                    str(file_count).zfill(8) + str(i).zfill(23) + "\n"
                    )
        file_count += 1
        print("file %s written" % out_filepath)
        rand_sleep_offs = randint(
            0,
            conf_data["cycle"]["random-sec-offset"]
            )
        sleep(conf_data["cycle"]["seconds"] + rand_sleep_offs)
        


if __name__ == "__main__":
    main()
