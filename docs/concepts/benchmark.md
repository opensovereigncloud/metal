# Benchmark Concept

`Benchmark` is a custom resource designed to provide valuable information about various tests that are executed on a server. It serves as a tool for understanding server performance and identifying any degradation that may be occurring.

With its Benchmark Controller and Benchmark CLI components, it enables the execution of multiple tests, provides resource management capabilities, and facilitates analysis of benchmark results. By leveraging the capabilities of Benchmark, users can effectively evaluate server performance, detect any degradation, and make informed decisions for optimization and improvement.

## Resources

1. Benchmark Controller:
    Its primary function is to monitor the `Benchmark` custom resource and perform calculations to determine the percentage difference between the old and new values for a given test. By analyzing these differences, the controller can provide insights into the performance of the server and detect any degradation that might be present.

2. Benchmark CLI:
    Is a command-line tool that is utilized to execute multiple tests with varying requirements. It is designed to facilitate the process of running tests and provides a convenient interface for interacting with the Kubernetes cluster. The Benchmark CLI leverages cgroups internally to separate processes and manage resource consumption efficiently.

    The Benchmark CLI offers the following features:

        1. Test Execution: The CLI allows users to specify and execute multiple tests with different configurations and requirements. This enables comprehensive benchmarking of the server's performance under various scenarios.

        2. Resource Management: The CLI utilizes cgroups to isolate processes and control resource allocation. This ensures that the tests are conducted in a controlled environment, preventing interference between different test executions and providing accurate performance metrics.

        3. Result Analysis: After the tests are completed, the CLI can generate reports of the benchmark results. This allows users to gain insights into the server's performance, identify any performance bottlenecks or degradation, and make informed decisions for optimization.

