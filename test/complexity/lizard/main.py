def calculate_grade(score):  # Cyclomatic Complexity = 4
    if score >= 90:
        return 'A'
    elif score >= 80:
        return 'B'
    elif score >= 70:
        return 'C'
    else:
        return 'F'

def is_valid_password(password):  # Cyclomatic Complexity = 8
    if len(password) < 8:
        return False
    
    has_upper = False
    has_lower = False
    has_digit = False
    
    for char in password:
        if char.isupper():
            has_upper = True
        elif char.islower():
            has_lower = True
        elif char.isdigit():
            has_digit = True
            
    return has_upper and has_lower and has_digit
