def inside(relative_dir, &block)
	Dir.chdir File.join(File.dirname(__FILE__), relative_dir), &block 
end

def run(cmd)
	system(cmd) || raise("#{cmd} failed")
end

task :default => [:build, :test]

task :build => :lexer do
	# run "go build"
end

task :lexer do
	inside "lexer/" do
		run "ragel -Z lexer.rl"
		run "go build"
	end
end

task :test => :build do
	inside 'lexer/' do
		run "go test"
	end
end