# Extended Brainfuck (BF) Specification

This document outlines the design and implementation details for an extended Brainfuck (BF) language with new features, alongside a roadmap for development.

## **New Features**

### **1. `=`: Program Return Value**
- **Functionality**: The `=` symbol uses the current byte on the tape as the program's return value.
- **Example**:
  ```brainfuck
  +++++   # Increment cell[0] by 5
  =       # Return 5 as the program output
  ```

### **2. Numeric Unrolling**
- **Functionality**: Numbers preceding `<`, `>`, `+`, `-` are unrolled into repeated operations.
- **Example**:
  ```brainfuck
  5<   # Translates to <<<<<
  10+  # Translates to ++++++++++
  ```

### **3. Absolute Pointer Positioning**
- **Functionality**: `%` sets the tape pointer to an absolute position.
- **Syntax**: `%(<absolute_position>)`
- **Example**:
  ```brainfuck
  %(1000)  # Moves pointer to position 1000
  ++++     # Increment cell[1000] by 4
  ```

### **4. External Function Declarations and Calls**
#### **Function Declaration Syntax**
- **Syntax**:
  ```
  FUNCTIONS
  function_name(offset1, offset2, ...) -> return_size;
  ENDFUNCTIONS
  ```
  - Each `offset` specifies the byte size of an input parameter.
  - `return_size` specifies the size of the return value in bytes.

- **Example**:
  ```
  FUNCTIONS
  print(1,1) -> 1;
  sum(1,1) -> 1;
  ENDFUNCTIONS
  ```

#### **Function Call Syntax**
- **Without Return**:
  ```brainfuck
  @(function_name, offset_sequence)
  ```
  - Calls the function but does not store its return value.

- **With Return**:
  ```brainfuck
  !@(function_name, offset_sequence)
  ```
  - Calls the function and stores the return value starting at the current pointer.

- **Offset Sequence**:
  - A list of relative offsets (e.g., `8+`, `5-`) or absolute positions using `%`.

- **Example**:
  ```brainfuck
  +++++      # Increment cell[0] by 5
  @(print, 8+, 5-)  # Call print with relative offsets 8 forward and 5 back
  !@(sum, %(1000), %(2000))  # Call sum with absolute offsets 1000 and 2000, store result at the current pointer
  ```

### **5. Tape Representation**
- **Size**: The tape is simulated as an array of `2^32` bytes.
- **Pointer**: The pointer can move freely within this range.

---

## **Roadmap**

### **1. Basic BF Syntax Working**
#### **Objective**:
- Implement core Brainfuck syntax: `>`, `<`, `+`, `-`, `[`, `]`, `,`, `.`.
- Add numeric unrolling for `<`, `>`, `+`, `-`.

#### **Milestones**:
1. Parse numeric prefixes and expand them into repeated operations.
2. Validate tape operations and handle boundary conditions.

### **2. Extending Syntax by Absolute Pointers**
#### **Objective**:
- Support absolute pointer syntax `%(<position>)`.
- Implement efficient tape access using absolute positioning.

#### **Milestones**:
1. Update parser to recognize `%(<position>)` syntax.
2. Modify runtime to process absolute pointer commands.
3. Add tests for edge cases (e.g., `%(-1)`, `%(2^32)`).

### **3. Implementing Function Calling**
#### **Objective**:
- Support function declarations and calls.
- Enable offset-based parameter passing and result storage.

#### **Milestones**:
1. Implement `FUNCTIONS` and `ENDFUNCTIONS` parsing.
2. Add syntax for function calls: `@(function_name, offset_sequence)` and `!@(function_name, offset_sequence)`.
3. Integrate runtime support for:
   - Parameter fetching based on offsets.
   - Return value storage.
4. Provide test cases for various function call scenarios.

---

## **Contributing**
Contributions are welcome! Please feel free to fork this repository and submit pull requests for bug fixes, new features, or documentation improvements.

