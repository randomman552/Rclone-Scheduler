#!/usr/bin/env python 
from requests import ConnectionError
from util import backup, get_parameter

def main():
    src = get_parameter("BACKUP_SOURCE", "/data")
    dest_remote = get_parameter("BACKUP_REMOTE", "remote")
    dest_path = get_parameter("BACKUP_DEST", "/backup")

    # Strip leading "/" causing malformed destination string
    # Interestingly this caused backups to not check for duplicates properly
    dest_path = dest_path.strip("/")

    dest = f"{dest_remote}/{dest_path}"

    operation = get_parameter("OPERATION", "backup")

    try:
        match operation:
            case "backup":
                response = backup(src, dest)
                json = response.json()
                jobId = json.get("jobid")
                print(f"Started backup job id: {jobId}")
            case _:
                raise NotImplementedError("Unrecognised operation")
    except ConnectionError as ex:
        print(f"Connection error: {ex}")
    

if __name__ == "__main__":
    main()