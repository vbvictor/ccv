// Function with cyclomatic complexity of 4
int processNumber(int num) {
  if (num < 0) {
    return -1;
  } else if (num == 0) {
    return 0;
  } else if (num > 100) {
    return 100;
  }
  return num;
}

// Function with cyclomatic complexity of 6
int validateAndProcessInput(int value) {
  if (value < -100) {
    return -100;
  } else if (value < 0) {
    return value * -1;
  } else if (value == 0) {
    return 0;
  } else if (value > 1000) {
    return 1000;
  } else if (value > 100) {
    return value / 2;
  }
  return value;
}
