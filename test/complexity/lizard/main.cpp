// Function with cyclomatic complexity of 4
// Complexity calculated as: 1 (base) + 3 (decision points)
int processNumber(int num) {
  if (num < 0) { // +1 complexity
    return -1;
  } else if (num == 0) { // +1 complexity
    return 0;
  } else if (num > 100) { // +1 complexity
    return 100;
  }
  return num;
}

// Function with cyclomatic complexity of 6
// Complexity calculated as: 1 (base) + 5 (decision points)
int validateAndProcessInput(int value) {
  if (value < -100) { // +1 complexity
    return -100;
  } else if (value < 0) { // +1 complexity
    return value * -1;
  } else if (value == 0) { // +1 complexity
    return 0;
  } else if (value > 1000) { // +1 complexity
    return 1000;
  } else if (value > 100) { // +1 complexity
    return value / 2;
  }
  return value;
}
