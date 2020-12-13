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
	seq_printf(m, "Memoria Total: %8lu MB\n",Convert(info.totalram)/1024);
	seq_printf(m, "Memoria Consumida: %8lu MB\n",(Convert(info.totalram)-Convert(info.freeram))*100);
	seq_printf(m, "Procentaje Consumo: %ld %%\n", (((Convert(info.totalram)-Convert(info.freeram))*100) / (Convert(info.totalram))*100)/100);
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

MODULE_LICENSE("GPL");
MODULE_DESCRIPTION("Modulo de Memoria - Sistemas Operativos 2");

static int __init memo_grupo22_init(void)
{
	printk(KERN_INFO "Hola mundo, somos el grupo 22(Los Cracks :v) y este es el monitor de memoria\n");
	proc_create("memo_grupo22", 0, NULL, &mem_info_fops);
	return 0;
}

static void __exit memo_grupo22_cleanup(void)
{
	remove_proc_entry("memo_grupo22", NULL);
	printk(KERN_INFO "Sayonara mundo, somos el grupo 22(Los Cracks :v) y este fue el monitor de memoria\n");
}

module_init(memo_grupo22_init);
module_exit(memo_grupo22_cleanup);
