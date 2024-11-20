def calculate_grade(score):  # Cyclomatic Complexity = 4
    # One path for initial entry
    if score >= 90:  # +1 for if condition
        return 'A'
    elif score >= 80:  # +1 for elif condition
        return 'B'
    elif score >= 70:  # +1 for elif condition
        return 'C'
    else:  # +1 for else condition
        return 'F'

def is_valid_password(password):  # Cyclomatic Complexity = 5
    # One path for initial entry
    if len(password) < 8:  # +1 for if condition
        return False
    
    has_upper = False
    has_lower = False
    has_digit = False
    
    for char in password:  # +1 for loop
        if char.isupper():  # +1 for if condition
            has_upper = True
        elif char.islower():  # +1 for elif condition
            has_lower = True
        elif char.isdigit():  # +1 for elif condition
            has_digit = True
            
    return has_upper and has_lower and has_digit
