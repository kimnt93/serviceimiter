import ctypes
from ctypes import c_char_p, c_int, POINTER, Structure

# Load the shared library
lib = ctypes.CDLL('./build/libratelimiter.so')

# Define C structs as Python ctypes structures
class RateLimitConfig(Structure):
    _fields_ = [
        ("AccountID", c_char_p),
        ("ServiceName", c_char_p),
        ("RequestPerSecond", c_int),
        ("RequestPerMinute", c_int),
        ("RequestPerHour", c_int),
        ("RequestPerDay", c_int),
        ("RequestPerWeek", c_int),
        ("RequestPerMonth", c_int),
        ("RequestPerYear", c_int),
    ]

class RateLimitRemaining(Structure):
    _fields_ = [
        ("AccountID", c_char_p),
        ("ServiceName", c_char_p),
        ("RequestPerSecond", c_int),
        ("RequestPerMinute", c_int),
        ("RequestPerHour", c_int),
        ("RequestPerDay", c_int),
        ("RequestPerWeek", c_int),
        ("RequestPerMonth", c_int),
        ("RequestPerYear", c_int),
    ]

class BucketConfig(Structure):
    _fields_ = [
        ("Key", c_char_p),
        ("PeriodType", c_int),
        ("Capacity", c_int),
        ("Ttl", c_int),
    ]

class RedisBucket(Structure):
    pass

class RateLimiter(Structure):
    pass

# Function prototypes
lib.NewRedisBucket.argtypes = [c_char_p, c_char_p, c_int]
lib.NewRedisBucket.restype = POINTER(RedisBucket)

lib.NewRateLimiter.argtypes = [POINTER(RedisBucket)]
lib.NewRateLimiter.restype = POINTER(RateLimiter)

lib.NewDefaultBucket.restype = POINTER(RedisBucket)

lib.IsAllow.argtypes = [POINTER(RateLimiter), POINTER(RateLimitConfig)]
lib.IsAllow.restype = c_int

lib.UpdateToken.argtypes = [POINTER(RateLimiter), POINTER(RateLimitConfig)]

lib.getErrorMessage.restype = c_char_p
lib.freeString.argtypes = [c_char_p]

# Wrapper functions
def new_redis_bucket(addr, password, db):
    return lib.NewRedisBucket(addr.encode('utf-8'), password.encode('utf-8'), db)

def new_rate_limiter(bucket):
    return lib.NewRateLimiter(bucket)

def new_default_bucket():
    return lib.NewDefaultBucket()

def is_allow(rate_limiter, config):
    return lib.IsAllow(rate_limiter, config)

def update_token(rate_limiter, config):
    lib.UpdateToken(rate_limiter, config)

def get_error_message():
    msg = lib.getErrorMessage()
    return ctypes.string_at(msg).decode('utf-8')

def free_string(s):
    lib.freeString(s.encode('utf-8'))

# Example usage
if __name__ == "__main__":
    bucket = new_default_bucket()
    rate_limiter = new_rate_limiter(bucket)
    config = RateLimitConfig(
        AccountID=b"myAccountID",
        ServiceName=b"myService",
        RequestPerSecond=10,
        RequestPerMinute=100,
        RequestPerHour=1000,
        RequestPerDay=10000,
        RequestPerWeek=70000,
        RequestPerMonth=300000,
        RequestPerYear=3600000
    )
    allowed = is_allow(rate_limiter, config)
    print(f"Allowed: {allowed}")
