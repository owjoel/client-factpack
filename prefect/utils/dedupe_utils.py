import pandas as pd
from pandas_dedupe import dedupe_dataframe
from dedupe import variables

# === CONFIG ===
INPUT_XLSX = "top_100_billionaires_expanded.xlsx"
SETTINGS_FILE = "dedupe_learned_settings"
TRAINING_FILE = "training.json"

FIELDS = [
    variables.String("Name"),
    variables.String("Companies"),
    variables.String("Primary Citizenship"),
]


def main():
    print("Loading data...")
    df = pd.read_excel(INPUT_XLSX)

    # Drop unused fields
    df = df.drop(columns=["Date of Birth"], errors="ignore")

    # Add identifier for training (optional)
    df["_is_existing"] = True

    print("Starting interactive labeling session (required for training)...")
    dedupe_dataframe(
        df,
        field_properties=FIELDS,
    )

    print("Training complete. Settings saved to:", SETTINGS_FILE)


if __name__ == "__main__":
    main()
