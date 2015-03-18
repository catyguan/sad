package bma.servicecall.core;

public class AppError extends RuntimeException {

	private static final long serialVersionUID = 1L;

	public AppError() {
		super();
	}

	public AppError(String message) {
		super(message);
	}

	public AppError(Throwable cause) {
		super(cause);
	}

	public AppError(String message, Throwable cause) {
		super(message, cause);
	}

	public static RuntimeException handle(Throwable t) {
		if (t instanceof RuntimeException) {
			RuntimeException ex = (RuntimeException) t;
			return ex;
		}
		return new AppError(t);
	}
}
