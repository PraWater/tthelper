# TimeTable Helper

![demo gif](./demo.gif)

## Usage

1. **Install locally**
   ```bash
   $ go install github.com/PraWater/tthelper@latest
   ```

2. **Populate the database using an excel file**

   Download the [excel file](./timetable.xlsx) from this repo. Last updated: *27th Dec*

   Default location for the Excel file is `/home/timetable.xlsx`.
   ```bash
   $ tthelper -refresh {path_to_excel_file}
   ```

3. **Create an input timetable as a TXT file**

     Example for 3-1 CS:
     ```txt
     CS F351 L1 T3
     CS F372 L1 T5
     CS F342 L1 P5
     CS F301 L1 T2
     ```

4. **Pass the TXT file to get the list of courses that do not clash with your timetable**

   Default location for the input timetable is `/home/input_tt.txt`.
   ```bash
   $ tthelper {path_to_input_tt}
   ```
