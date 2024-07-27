# TimeTable Helper

![demo gif](./demo.gif)

## Usage

1. **Install Locally**
   ```sh
   $ go install github.com/PraWater/tthelper@latest
   ```

2. **Populate the Database Using an Excel File**

   Download the [excel file](./timetable.xlsx) from this repo. Last updated: *27th July*

   Default location for the Excel file is `/home/timetable.xlsx`.
   ```sh
   $ tthelper -refresh {path_to_excel_file}
   ```

3. **Create an Input Timetable as a TXT File**

     Example for 3-1 CS:
     ```txt
     CS F351 L1 T3
     CS F372 L1 T5
     CS F342 L1 P5
     CS F301 L1 T2
     ```

4. **Pass the TXT File to Get the List of Courses that do not clash with your Timetable**

   Default location for the input timetable is `/home/input_tt.txt`.
   ```sh
   $ tthelper {path_to_input_tt}
   ```
