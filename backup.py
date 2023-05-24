from methods import backup, get_parameter

def main():
    src = get_parameter("BACKUP_SOURCE", "/data")
    dest_remote = get_parameter("BACKUP_REMOTE", "remote")
    dest_path = get_parameter("BACKUP_DEST", "/backup")

    dest = f"{dest_remote}/{dest_path}"

    operation = get_parameter("OPERATION", "backup")

    match operation:
        case "backup":
            response = backup(src, dest)
            json = response.json()
            jobId = json.get("jobid")
            print(f"Backup job started with id: {jobId}")
        case _:
            raise NotImplementedError("Unrecognised operation")
    

if __name__ == "__main__":
    main()