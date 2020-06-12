#include <linux/init.h>
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/fs.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <asm/uaccess.h>
#include <linux/hugetlb.h>
#include <linux/mm.h>
#include <linux/mman.h>
#include <linux/mmzone.h>
#include <linux/syscalls.h>
#include <linux/swap.h>
#include <linux/swapfile.h>
#include <linux/vmstat.h>
#include <linux/atomic.h>


struct sysinfo info;

static int leer_memoria(struct seq_file *m, void *v){

    #define Convert(x) ((x) << (PAGE_SHIFT - 10))
	si_meminfo(&info); 
	seq_printf(m, "Carnet1: 201122826\n");
	seq_printf(m, "Nombre1: Ivan Alfonso Lopez Medina\n");
    seq_printf(m, "Carnet2: 200312755\n");
	seq_printf(m, "Nombre2: Herminio Rolando García Sánchez\n");	
	seq_printf(m, "Memoria Total: %8lu MB\n",Convert(info.totalram)/1024);
	seq_printf(m, "Memoria Libre: %8lu kB\n",Convert(info.freeram)/1024);
	seq_printf(m, "Memoria Usada: %ld %%\n", (((Convert(info.totalram)-Convert(info.freeram))*100) / (Convert(info.totalram))*100)/100);
	#undef K
	return 0;

}

static int mem_info_open(struct inode *inode, struct file *file){
	return single_open(file, leer_memoria, NULL);
}

static const struct file_operations mem_info_fops = {
	.owner = THIS_MODULE,
	.open = mem_info_open,
	.read = seq_read,
	.llseek = seq_lseek,
	.release = single_release,
};

static int __init memo_201122826_init(void)
{
	printk(KERN_INFO "201122826 - 200312755\n");

	proc_create("memo_201122826", 0, NULL, &mem_info_fops);
	return 0;
}

static void __exit memo_201122826_cleanup(void)
{
	remove_proc_entry("memo_201122826", NULL);
	printk(KERN_INFO "Sistemas Operativos 1\n");
}

module_init(memo_201122826_init);
module_exit(memo_201122826_cleanup);
