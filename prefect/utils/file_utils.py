import os
import csv
import json
import yaml


def save_to_csv(data, folder, file_name):
    """Save a list of dictionaries to a CSV file in the specified folder."""
    if not data:
        print("No data to save.")
        return None

    os.makedirs(os.path.join("data", folder), exist_ok=True)
    filename = os.path.join("data", folder, f"{file_name}.csv")

    with open(filename, "w", newline="", encoding="utf-8") as f:
        writer = csv.DictWriter(f, fieldnames=list(data[0].keys()))
        writer.writeheader()
        writer.writerows(data)

    print(f"Saved {len(data)} individuals to {filename}")
    return filename


def save_to_text(data, folder_name, file_name):
    """Save a string to a text file in the specified folder."""
    os.makedirs(folder_name, exist_ok=True)
    file_path = os.path.join(folder_name, file_name)

    with open(file_path, "w", encoding="utf-8") as f:
        f.write(data)

    return file_path


def save_to_json(data, folder_name, file_name):
    """Save a dictionary or list to a JSON file in the specified folder."""
    os.makedirs(folder_name, exist_ok=True)
    file_path = os.path.join(folder_name, file_name)

    with open(file_path, "w", encoding="utf-8") as f:
        json.dump(data, f, indent=2)

    return file_path


def read_text_file(file_path):
    """Read and return contents of a UTF-8 encoded text file."""
    with open(file_path, "r", encoding="utf-8") as file:
        return file.read().strip()


def read_yaml_file(file_path):
    """Load and return YAML content from the given file."""
    with open(file_path, "r", encoding="utf-8") as file:
        return yaml.safe_load(file)


def read_json_file(file_path):
    """Load and return JSON content from the given file."""
    with open(file_path, "r", encoding="utf-8") as file:
        return json.load(file)
